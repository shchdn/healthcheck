package healthcheck

import (
	"github.com/go-playground/assert/v2"
	"net/http"
	"template/internal/common"
	"testing"
	"time"
)

func TestGetMin(t *testing.T) {
	okCode := 200
	minLatency := time.Second * 1
	cases := []common.SiteInfo{
		{
			Url:     "foo",
			Status:  &okCode,
			Latency: minLatency,
			Error:   nil,
		},
		{
			Url:     "bar",
			Status:  &okCode,
			Latency: time.Second * 5,
			Error:   nil,
		},
		{
			Url:     "avg",
			Status:  &okCode,
			Latency: time.Second * 2,
			Error:   nil,
		},
	}
	s := Service{}
	s.stats = initHealthcheckStatistics()
	s.client = &http.Client{}
	data := make(map[url]common.SiteInfo)
	for _, v := range cases {
		data[v.Url] = v
	}
	s.stats.UpdateStatistics(data)
	actual, err := s.stats.GetMin()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, minLatency, actual.Latency)
}

func TestGetMax(t *testing.T) {
	okCode := 200
	maxLatency := time.Second * 5
	cases := []common.SiteInfo{
		{
			Url:     "foo",
			Status:  &okCode,
			Latency: time.Second * 1,
			Error:   nil,
		},
		{
			Url:     "bar",
			Status:  &okCode,
			Latency: maxLatency,
			Error:   nil,
		},
		{
			Url:     "avg",
			Status:  &okCode,
			Latency: time.Second * 2,
			Error:   nil,
		},
	}
	s := Service{}
	s.stats = initHealthcheckStatistics()
	s.client = &http.Client{}
	data := make(map[url]common.SiteInfo)
	for _, v := range cases {
		data[v.Url] = v
	}
	s.stats.UpdateStatistics(data)
	actual, err := s.stats.GetMax()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, maxLatency, actual.Latency)
}
