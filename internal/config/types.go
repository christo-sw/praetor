package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Ping      Ping      `json:"ping"`
	Speedtest Speedtest `json:"speedtest"`
}

type Ping struct {
	Targets []PingTarget `json:"targets"`
}

type PingTarget struct {
	Endpoint   string `json:"endpoint"`
	IntervalMS int    `json:"intervalMs"`
}

type Speedtest struct {
	Targets    []SpeedtestTarget `json:"targets"`
	Unit       string            `json:"unit"`
	IntervalMS int               `json:"intervalMs"`
}

type SpeedtestTarget struct {
	ServerID int `json:"serverID"`
}

func ParseConfig() (*Config, error) {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}
