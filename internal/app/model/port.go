package model

type Port struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias,omitempty"`
	Regions     []string  `json:"regions,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs,omitempty"`
	Code        string    `json:"code"`
}
