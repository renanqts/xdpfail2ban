package config

// Urlregexp struct
type Urlregexp struct {
	Regexp string `yaml:"regexp"`
}

// Rules struct
type Rules struct {
	Bantime    string      `yaml:"bantime"`
	Findtime   string      `yaml:"findtime"`
	Maxretry   int         `yaml:"maxretry"`
	Urlregexps []Urlregexp `yaml:"urlregexps"`
}

// Config struct
type Config struct {
	Rules         Rules  `yaml:"port"`
	LogLevel      string `yaml:"loglevel"`
	XDPDropperURL string `yaml:"xdpdropperurl"`
}
