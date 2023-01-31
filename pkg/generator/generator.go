package generator

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/crodriguezde/rtdashs/pkg/events"
)

// Generate events for testing purposes
func RunCpuGenerator(ctx context.Context, send chan *events.Cpu) {
	log.Printf("Starting cpu generator")
	ticker := time.NewTicker(1 * time.Second)
	go func() {
	outter:
		for {
			select {
			case <-ticker.C:
				log.Printf("Tick!")
				send <- events.NewCpu(rand.Intn(100))
			case <-ctx.Done():
				ticker.Stop()
				break outter
			}
		}
	}()

}
