package sink

import (
	"context"
	"log"
	"time"

	kafkaPlayloads "github.com/crodriguezde/rtdashs/pkg/kafkaPayloads"
)

type Server struct {
	recv    chan *kafkaPlayloads.Cpu
	send    chan map[string]int64
	ctx     context.Context
	cpudata Metric
}

func NewServer(ctx context.Context, recv chan *kafkaPlayloads.Cpu) *Server {
	return &Server{
		recv:    recv,
		ctx:     ctx,
		cpudata: *NewMetric(ctx),
		send:    make(chan map[string]int64),
	}
}

func (s *Server) Start() {
	go func() {
		log.Print("starting sink")
		ticker := time.NewTicker(10 * time.Second)
	outter:
		for {
			select {
			case <-ticker.C:
				s.cpudata.Aggregate()
				// dont send data if avg array is empty
				if len(s.cpudata.avg) > 0 {
					s.send <- s.cpudata.avg
				}
			case msg := <-s.recv:
				s.cpudata.Add(msg.Host, msg.Timestamp)
			case <-s.ctx.Done():
				break outter
			}
		}
	}()
}

func (s *Server) GetChannel() chan map[string]int64 {
	return s.send
}
