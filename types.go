package dvb

type Status struct {
	Code    string `json:"Code"`
	Message string `json:"Message,omitempty"`
}

type Diva struct {
	Number  string `json:"Number"`
	Network string `json:"Network"`
}

type Platform struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
}
