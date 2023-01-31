package kafkaPlayloads

type Generic struct {
	Value string `json:"value"`
	Key   string `json:"string"`
}

type Cpu struct {
	UsageActive  float32 `json:"usage_active"`
	UsageIdle    float32 `json:"usage_idle"`
	CoreId       string  `json:"core_id"`
	Cpu          string  `json:"cpu"`
	Environment  string  `json:"environment"`
	Host         string  `json:"host"`
	Organization string  `json:"organization"`
	PhysicalId   string  `json:"physical_id"`
	Timestamp    float32 `json:"timestamp"`
}
