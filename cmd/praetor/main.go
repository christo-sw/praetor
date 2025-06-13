package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/christo-sw/praetor/internal/config"
	"github.com/christo-sw/praetor/internal/speedtest"
	probing "github.com/prometheus-community/pro-bing"
)

type Ping struct {
	Seq            int
	SentPacket     *probing.Packet
	ReceivedPacket *probing.Packet
	Stats          *probing.Statistics
}

func (p Ping) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Seq: %v\n", p.Seq))

	if p.SentPacket != nil {
		sb.WriteString(fmt.Sprintf("Sent Packet ID: %v\n", p.SentPacket.ID))
	}

	if p.ReceivedPacket != nil {
		sb.WriteString(fmt.Sprintf("Received Packet ID: %v, RTT: %v\n", p.SentPacket.ID, p.ReceivedPacket.Rtt))
	}

	if p.Stats != nil {
		sb.WriteString(fmt.Sprintf("Stats: Round-trip Min/Avg/Max/StdDev = %v/%v/%v/%v\n",
			p.Stats.MinRtt, p.Stats.AvgRtt, p.Stats.MaxRtt, p.Stats.StdDevRtt))
	}

	return sb.String()
}

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	pingers := make([]*probing.Pinger, 0)

	err = registerPingers(&pingers, config)
	if err != nil {
		panic(err)
	}

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			stopPingers(&pingers)
		}
	}()

	go speedtestThread(config)
	runPingers(&pingers)

	// pingResponses := make([]Ping, 0, 2048)
}

func registerPingers(pingers *[]*probing.Pinger, cfg *config.Config) error {
	for _, target := range cfg.Ping.Targets {
		pinger, err := probing.NewPinger(target.Endpoint)
		if err != nil {
			return fmt.Errorf("failed to create pinger: %v", err)
		}

		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		pinger.Interval = time.Duration(target.IntervalMS) * time.Millisecond

		pinger.OnSend = func(pkt *probing.Packet) {
			fmt.Printf("Seq: %v, IP: %v\n", pkt.Seq, pkt.Addr)
		}

		pinger.OnRecv = func(pkt *probing.Packet) {
			fmt.Printf("Seq: %v, IP: %v, RTT: %v\n", pkt.Seq, pkt.Addr, pkt.Rtt)
		}

		pinger.OnFinish = func(stats *probing.Statistics) {
		}

		*pingers = append(*pingers, pinger)
	}

	return nil
}

func runPingers(pingers *[]*probing.Pinger) {
	wg := sync.WaitGroup{}
	wg.Add(len(*pingers))

	for _, pinger := range *pingers {
		go func() {
			_ = pinger.Run()
			wg.Done()
		}()
	}

	wg.Wait()
}

func stopPingers(pingers *[]*probing.Pinger) {
	for _, pinger := range *pingers {
		pinger.Stop()
	}
}

func speedtestThread(cfg *config.Config) {
	unit := speedtest.ParseUnit(cfg.Speedtest.Unit)

	for {
		for _, target := range cfg.Speedtest.Targets {
			speedtestResults, err := speedtest.RunSpeedtest(target.ServerID)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Download: %v Mbps, Upload: %v Mbps\n", speedtestResults.DownloadSpeed(unit), speedtestResults.UploadSpeed(unit))
		}

		time.Sleep(time.Duration(cfg.Speedtest.IntervalMS) * time.Millisecond)
	}
}
