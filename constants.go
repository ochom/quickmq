package main

import (
	"os"
	"strconv"
	"time"
)

var (
	// CronJobInterval is the interval at which the cron job runs
	CronJobInterval = getCronInterval(1 * time.Second)
)

func getCronInterval(d time.Duration) time.Duration {
	val, ok := os.LookupEnv("CRON_INTERVAL")
	if !ok {
		return d
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return d
	}

	return time.Duration(num) * time.Minute
}
