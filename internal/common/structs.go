package common

import "time"

type SiteInfo struct {
	Url     string
	Status  *int
	Latency time.Duration
	Error   error
}
