package publisher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/cloudzxx/crypto-match-core/internal/config"
)

type Publisher struct {
	writers map[string]*kafka.Writer
	topics  config.KafkaTopics
}

func NewPublisher(cfg config.KafkaConfig) *Publisher {
	p := &Publisher{
		writers: make(map[string]*kafka.Writer),
		topics:  cfg.Topics,
	}

	for name, topic := range map[string]string{
		"orders":   cfg.Topics.Orders,
		"deals":    cfg.Topics.Deals,
		"balances": cfg.Topics.Balances,
	} {
		p.writers[name] = &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: time.Duration(cfg.BatchMs) * time.Millisecond,
			Async:        true,
		}
	}

	return p
}

type OrderMessage struct {
	Type    string `json:"type"`
	OrderID uint64 `json:"order_id"`
	UserID  uint32 `json:"user_id"`
	Market  string `json:"market"`
	Side    string `json:"side"`
	Price   string `json:"price"`
	Amount  string `json:"amount"`
	Left    string `json:"left"`
}

type DealMessage struct {
	TakerID  uint64 `json:"taker_id"`
	MakerID  uint64 `json:"maker_id"`
	Market   string `json:"market"`
	Price    string `json:"price"`
	Amount   string `json:"amount"`
	MakerFee string `json:"maker_fee"`
	TakerFee string `json:"taker_fee"`
}

type BalanceMessage struct {
	UserID    uint32 `json:"user_id"`
	Asset     string `json:"asset"`
	Available string `json:"available"`
	Freeze    string `json:"freeze"`
}

func (p *Publisher) PublishOrder(ctx context.Context, msg *OrderMessage) error {
	data, _ := json.Marshal(msg)
	return p.writers["orders"].WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Market),
		Value: data,
	})
}

func (p *Publisher) PublishDeal(ctx context.Context, msg *DealMessage) error {
	data, _ := json.Marshal(msg)
	return p.writers["deals"].WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Market),
		Value: data,
	})
}

func (p *Publisher) PublishBalance(ctx context.Context, msg *BalanceMessage) error {
	data, _ := json.Marshal(msg)
	return p.writers["balances"].WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Asset),
		Value: data,
	})
}

func (p *Publisher) Close() error {
	for _, w := range p.writers {
		w.Close()
	}
	return nil
}