package healthcheck

import "sync"

type endpointName = string
type visitedCounter = int

type RequestStatistics struct {
	mutex *sync.RWMutex
	stats map[endpointName]visitedCounter
}

func initRequestStatistics() *RequestStatistics {
	stats := make(map[endpointName]visitedCounter)
	return &RequestStatistics{
		mutex: &sync.RWMutex{},
		stats: stats,
	}
}

func (r *RequestStatistics) IncrementStats(endpointName string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.stats[endpointName]++
}

func (r *RequestStatistics) GetStats() map[endpointName]visitedCounter {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.stats
}
