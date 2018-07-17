package config

type Config struct {
	Desc      string
	Listen    string
	Bootstrap bool
	NodeID    string
}

func NewConfig(listen, nodeid string, bootstrap bool) *Config {
	return &Config{
		Listen:    listen,
		NodeID:    nodeid,
		Bootstrap: bootstrap,
		Desc:      "",
	}
}
