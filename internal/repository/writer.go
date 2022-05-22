package repository

import (
	"context"
	"encoding/json"
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
	select {
	case w.logs <- opt:
	default:
	}
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

	for _, log := range w.batch {
		data, _ := json.Marshal(log.Data)
		w.pg.Exec(ctx, `
			INSERT INTO operlog (opt_type, opt_data, seq_id)
			VALUES ($1, $2, $3)
		`, log.Type, data, w.seqID)
		w.seqID++
	}

	w.batch = w.batch[:0]
}

func (w *Writer) Stop() {
	close(w.done)
	time.Sleep(time.Second)
}

func (w *Writer) NextSeqID() int64 {
	return w.seqID
}
