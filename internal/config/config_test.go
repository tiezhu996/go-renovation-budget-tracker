package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("RENO_WORKERS", "")
	t.Setenv("RENO_BATCH_SIZE", "")
	t.Setenv("RENO_ALERT_THRESHOLD", "")
	c := Load()
	if c.Workers != 2 || c.BatchSize != 2 || c.Threshold != 0.9 {
		t.Fatalf("bad config %+v", c)
	}
}
