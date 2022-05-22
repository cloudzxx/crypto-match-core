package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Grpc   GrpcConfig   `yaml:"grpc"`
	PG     PGConfig     `yaml:"pg"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	Slice  SliceConfig  `yaml:"slice"`
	Assets []AssetConfig `yaml:"assets"`
	Markets []MarketConfig `yaml:"markets"`
}

type GrpcConfig struct {
	Addr string `yaml:"addr"`
}

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	MaxConns int    `yaml:"max_conns"`
}

type KafkaConfig struct {
	Brokers  []string    `yaml:"brokers"`
	Topics   KafkaTopics `yaml:"topics"`
	BatchMs  int         `yaml:"batch_ms"`
}

type KafkaTopics struct {
	Orders    string `yaml:"orders"`
	Deals     string `yaml:"deals"`
	Balances  string `yaml:"balances"`
}

type SliceConfig struct {
	Interval int `yaml:"interval"`
	KeepTime  int `yaml:"keeptime"`
}

type AssetConfig struct {
	Name     string `yaml:"name"`
	PrecSave int    `yaml:"prec_save"`
	PrecShow int    `yaml:"prec_show"`
}

type MarketConfig struct {
	Name       string `yaml:"name"`
	Stock      string `yaml:"stock"`
	Money      string `yaml:"money"`
	StockPrec  int    `yaml:"stock_prec"`
	MoneyPrec  int    `yaml:"money_prec"`
	MinAmount  string `yaml:"min_amount"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}