package main

import (
	"fmt"
	"log"

	"github.com/cloudzxx/crypto-match-core/internal/config"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	fmt.Println("Crypto Match Core")
	fmt.Printf("gRPC addr: %s\n", cfg.Grpc.Addr)
	fmt.Printf("PG host: %s:%d/%s\n", cfg.PG.Host, cfg.PG.Port, cfg.PG.Name)
	fmt.Printf("Kafka brokers: %v\n", cfg.Kafka.Brokers)
	fmt.Printf("Markets: %d configured\n", len(cfg.Markets))
}