package events

type Cpu struct {
	Value int `json:"value"`
}

func NewCpu(v int) *Cpu {
	return &Cpu{Value: v}
}
