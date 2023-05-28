package healthcheck

import (
	"errors"
	"net/http"
	"sync"
	"template/internal/common"
)

type url = string

type Statistics struct {
	lastData map[url]common.SiteInfo
	min      *common.SiteInfo
	max      *common.SiteInfo
	mutex    *sync.RWMutex
}

func initHealthcheckStatistics() *Statistics {
	return &Statistics{
		lastData: make(map[url]common.SiteInfo),
		min:      nil,
		max:      nil,
		mutex:    &sync.RWMutex{},
	}
}

func (s *Statistics) UpdateStatistics(data map[url]common.SiteInfo) {
	var minLatencySite, maxLatencySite common.SiteInfo
	for _, info := range data {
		if info.Error != nil || *info.Status != http.StatusOK {
			continue
		}
		if minLatencySite.Latency == 0 || info.Latency < minLatencySite.Latency {
			minLatencySite = info
		}
		if maxLatencySite.Latency == 0 || info.Latency > maxLatencySite.Latency {
			maxLatencySite = info
		}
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastData = data
	s.min = &minLatencySite
	s.max = &maxLatencySite
}

func (s *Statistics) GetSiteInfo(url string) (*common.SiteInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.lastData[url]; ok {
		if v.Error != nil {
			return nil, errors.New("no data")
		}
		return &v, nil
	}
	return nil, errors.New("unknown site")
}

func (s *Statistics) GetMin() (*common.SiteInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.min != nil {
		return s.min, nil
	}
	return nil, errors.New("min is not ready")
}

func (s *Statistics) GetMax() (*common.SiteInfo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.max != nil {
		return s.max, nil
	}
	return nil, errors.New("max is not ready")
}
