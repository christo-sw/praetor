package speedtest

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

func RunSpeedtest(serverID int) (*SpeedtestResult, error) {
	var cmd *exec.Cmd
	if serverID != 0 {
		cmd = exec.Command("speedtest", "-s", strconv.FormatInt(int64(serverID), 10), "-f", "json")
	} else {
		cmd = exec.Command("speedtest", "-f", "json")
	}

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	jsonString := string(data)
	jsonLines := strings.Split(jsonString, "\n")

	var result SpeedtestResult
	for _, s := range jsonLines {
		err = json.Unmarshal([]byte(s), &result)
		if err != nil {
			continue
		}
	}

	return &result, nil
}

func (sr *SpeedtestResult) DownloadSpeed(u SpeedUnit) float64 {
	return getSpeed(sr.Download.Bytes, sr.Download.Elapsed, u)
}

func (sr *SpeedtestResult) UploadSpeed(u SpeedUnit) float64 {
	return getSpeed(sr.Upload.Bytes, sr.Upload.Elapsed, u)
}

func getSpeed(bytes, elapsed int64, u SpeedUnit) float64 {
	var byps float64 = float64(bytes) / (float64(elapsed) / 1000)

	switch u {
	case Bps:
		return byps * 8
	case Byps:
		return byps
	case Kbps:
		return (byps * 8) / 1_000
	case KBps:
		return byps / 1_000
	case Mbps:
		return (byps * 8) / 1_000_000
	case MBps:
		return byps / 1_000_000
	case Gbps:
		return (byps * 8) / 1_000_000_000
	case GBps:
		return byps / 1_000_000_000
	}

	return 0
}
