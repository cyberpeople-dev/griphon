package internal

type Endpoint struct {
	Path string `yaml:"path"`
}

type Payload struct {
	Name string                 `yaml:"name"`
	Data map[string]interface{} `yaml:"data"`
}

type Config struct {
	UserAgents []string   `yaml:"user_agents"`
	Endpoints  []Endpoint `yaml:"endpoints"`
	Payloads   []Payload  `yaml:"payloads"`
}

type ResponseResult struct {
	URL    string
	Method string
	Status int
	Body   string
} 