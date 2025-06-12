package speedtest

import (
	"time"
)

type SpeedUnit int

const (
	Bps SpeedUnit = iota
	Byps
	Kbps
	KBps
	Mbps
	MBps
	Gbps
	GBps
)

type SpeedtestLog struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

type SpeedtestResult struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
		Low     float64 `json:"low"`
		High    float64 `json:"high"`
	} `json:"ping"`
	Download struct {
		Bandwidth int64 `json:"bandwidth"`
		Bytes     int64 `json:"bytes"`
		Elapsed   int64 `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"download"`
	Upload struct {
		Bandwidth int64 `json:"bandwidth"`
		Bytes     int64 `json:"bytes"`
		Elapsed   int64 `json:"elapsed"`
		Latency   struct {
			Iqm    float64 `json:"iqm"`
			Low    float64 `json:"low"`
			High   float64 `json:"high"`
			Jitter float64 `json:"jitter"`
		} `json:"latency"`
	} `json:"upload"`
	PacketLoss int64  `json:"packetLoss"`
	Isp        string `json:"isp"`
	Interface  struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		ID       int64  `json:"id"`
		Host     string `json:"host"`
		Port     int64  `json:"port"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		IP       string `json:"ip"`
	} `json:"server"`
	Result struct {
		ID        string `json:"id"`
		URL       string `json:"url"`
		Persisted bool   `json:"persisted"`
	} `json:"result"`
}
