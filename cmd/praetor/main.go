package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

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
	pinger, err := probing.NewPinger("8.8.8.8")
	if err != nil {
		panic(err)
	}

	go SpeedtestThread(30 * time.Second)

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()

	pingResponses := make([]Ping, 0, 2048)

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Interval = 100 * time.Millisecond
	m := sync.Mutex{}

	pinger.OnSend = func(pkt *probing.Packet) {
		m.Lock()
		pingResponses = append(pingResponses, Ping{
			Seq:        pkt.Seq,
			SentPacket: pkt,
		})
		m.Unlock()
	}

	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("RTT: %v\n", pkt.Rtt)
		m.Lock()
		for i := len(pingResponses) - 1; i >= 0; i++ {
			if pingResponses[i].Seq == pkt.Seq {
				pingResponses[i].ReceivedPacket = pkt
				pingResponses[i].Stats = pinger.Statistics()
				break
			}
		}
		m.Unlock()
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("DONE!\n")
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		panic(err)
	}
}

func SpeedtestThread(d time.Duration) {
	for {
		speedtestResults, err := speedtest.RunSpeedtest()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Download: %v Mbps, Upload: %v Mbps\n", speedtestResults.DownloadSpeed(speedtest.Mbps), speedtestResults.UploadSpeed(speedtest.Mbps))

		time.Sleep(d)
	}
}
