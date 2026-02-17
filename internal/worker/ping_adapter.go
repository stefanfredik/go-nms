package worker

import (
	"os/exec"
	"time"
)

// PingAdapter encapsulates ping logic
type PingAdapter struct{}

func (p *PingAdapter) Ping(ip string) (time.Duration, bool) {
	start := time.Now()
	// ping -c 1 -W 1 <ip>
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()
	
	elapsed := time.Since(start)
	if err != nil {
		return 0, false
	}
	return elapsed, true
}
