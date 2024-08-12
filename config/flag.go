package config

import "flag"

func parseFlag(c *Config) {
	flag.StringVar(&c.Addr, "a", "localhost:8080", "Net address host:port")
	flag.StringVar(&c.DatabaseURI, "d", "", "db path")

	flag.Parse()
}
