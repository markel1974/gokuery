package config

type Config struct {
	EscapeQueryString bool   `json:"escapeQueryString"`
	DateFormatTZ      string `json:"dateFormatTZ"`
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg * Config) HasTimeZone() bool {
	return len(cfg.DateFormatTZ) > 0
}

func (cfg *Config) GetTimeZone() string {
	//t := time.Now()
	//zone, offset := t.Zone()
	//fmt.Println(zone, offset)
	return cfg.DateFormatTZ
}