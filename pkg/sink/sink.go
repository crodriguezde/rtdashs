package sink

import (
	"context"
	"fmt"
	"log"
	"time"

	kafkaPlayloads "github.com/crodriguezde/rtdashs/pkg/kafkaPayloads"
)

type Server struct {
	recv chan *kafkaPlayloads.Cpu
	// Map of channels, one per session
	sessions map[string]chan map[string]int64
	ctx      context.Context
	cpudata  Metric
}

func NewServer(ctx context.Context, recv chan *kafkaPlayloads.Cpu) *Server {
	return &Server{
		recv:     recv,
		ctx:      ctx,
		cpudata:  *NewMetric(ctx),
		sessions: make(map[string]chan map[string]int64),
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
				log.Printf("Tick aggregating")
				s.cpudata.Aggregate()
				for id, session := range s.sessions {
					log.Printf("Sending data to %s", id)
					session <- s.cpudata.avg
				}
			case msg := <-s.recv:
				s.cpudata.Add(msg.Host, msg.Timestamp)
			case <-s.ctx.Done():
				break outter
			}
		}
	}()
}

func (s *Server) Register(id string) (chan map[string]int64, error) {
	var ch chan map[string]int64
	if _, ok := s.sessions[id]; !ok {
		ch = make(chan map[string]int64)
		log.Printf("Stream with id %s registered", id)
		s.sessions[id] = ch
		return ch, nil
	}

	return ch, fmt.Errorf("id already in session")
}

func (s *Server) Deregister(id string) {
	delete(s.sessions, id)
	log.Printf("Stream with id %s deregistered", id)
}
