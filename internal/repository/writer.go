package repository

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
)

type OperType string

const (
	OperTypeOrderPut     OperType = "ORDER_PUT"
	OperTypeOrderUpdate  OperType = "ORDER_UPDATE"
	OperTypeOrderFinish  OperType = "ORDER_FINISH"
	OperTypeBalanceChange OperType = "BALANCE_CHANGE"
)

type OperLog struct {
	Type   OperType
	Data   interface{}
	SeqID  int64
}

type Writer struct {
	pg       *PG
	logs     chan *OperLog
	batch    []*OperLog
	batchMs  int
	mu       sync.Mutex
	done     chan struct{}
	seqID    int64
}

func NewWriter(pg *PG, batchMs int) *Writer {
	w := &Writer{
		pg:      pg,
		logs:    make(chan *OperLog, 10000),
		batch:   make([]*OperLog, 0, 100),
		batchMs: batchMs,
		done:    make(chan struct{}),
	}

	go w.run()

	return w
}

func (w *Writer) Write(opt *OperLog) {
	w.logs <- opt
}

func (w *Writer) run() {
	ticker := time.NewTicker(time.Duration(w.batchMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case log := <-w.logs:
			w.batch = append(w.batch, log)
			if len(w.batch) >= 100 {
				w.flush()
			}
		case <-ticker.C:
			w.flush()
		case <-w.done:
			w.flush()
			return
		}
	}
}

func (w *Writer) flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.batch) == 0 {
		return
	}

	ctx := context.Background()

	for _, l := range w.batch {
		data, err := json.Marshal(l.Data)
		if err != nil {
			log.Printf("operlog marshal error: %v", err)
			continue
		}
		if err := w.pg.Exec(ctx, `
			INSERT INTO operlog (opt_type, opt_data, seq_id)
			VALUES ($1, $2, $3)
		`, l.Type, data, w.seqID); err != nil {
			log.Printf("operlog write error: %v", err)
		}
		w.seqID++
	}

	w.batch = w.batch[:0]
}

func (w *Writer) Stop() {
	close(w.done)
	// 等待 run goroutine 处理完剩余日志
	done := make(chan struct{})
	go func() {
		for len(w.logs) > 0 || len(w.batch) > 0 {
			time.Sleep(10 * time.Millisecond)
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		log.Printf("operlog drain timeout, %d logs may be lost", len(w.logs)+len(w.batch))
	}
}

func (w *Writer) NextSeqID() int64 {
	return w.seqID
}
