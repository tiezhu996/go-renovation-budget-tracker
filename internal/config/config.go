package config

import (
	"os"
	"strconv"
)

type Config struct {
	Workers   int
	BatchSize int
	Threshold float64
}

func Load() *Config {
	return &Config{
		Workers:   getInt("RENO_WORKERS", 2),
		BatchSize: getInt("RENO_BATCH_SIZE", 2),
		Threshold: getFloat("RENO_ALERT_THRESHOLD", 0.9),
	}
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func getFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f <= 0 {
		return def
	}
	return f
}
