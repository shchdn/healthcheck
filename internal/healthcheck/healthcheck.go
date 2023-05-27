package healthcheck

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"sync"
	"template/internal/common"
	"time"
)

const ParseDelaySec = 60
const RequestDelaySec = 10

type SiteList []string

type Service struct {
	stats        *Statistics
	requestStats *RequestStatistics
	client       *http.Client
	siteList     SiteList
}

func New() *Service {
	p := Service{}
	p.stats = initHealthcheckStatistics()
	p.client = &http.Client{}
	p.siteList = p.getSiteList()
	p.requestStats = initRequestStatistics()
	return &p
}

func (s *Service) getSiteList() SiteList {
	var siteList SiteList
	jsonList, err := os.Open("./cmd/http/site-list.json")
	if err != nil {
		panic("cant open site list")
	}
	byteValue, err := io.ReadAll(jsonList)
	if err != nil {
		panic("cant read site list")
	}
	err = jsonList.Close()
	if err != nil {
		panic("close failed")
	}
	err = json.Unmarshal(byteValue, &siteList)
	if err != nil {
		panic("cant unmarshal")
	}
	return siteList
}

func (s *Service) Start() {
	channel := make(chan common.SiteInfo, len(s.siteList))
	for {
		wg := sync.WaitGroup{}
		for _, siteUrl := range s.siteList {
			wg.Add(1)
			go s.handleUrl(siteUrl, &wg, channel)
		}
		wg.Wait()
		lastData := make(map[url]common.SiteInfo)
		for i := 0; i < len(s.siteList); i++ {
			v := <-channel
			lastData[v.Url] = v
		}
		s.stats.UpdateStatistics(lastData)
		time.Sleep(ParseDelaySec * time.Second)
	}
}

func (s *Service) handleUrl(siteUrl string, wg *sync.WaitGroup, ch chan common.SiteInfo) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*RequestDelaySec))
	defer cancel()
	siteInfo, err := s.requestUrl(ctx, siteUrl)
	if err != nil {
		logrus.Error(err)
	}
	ch <- *siteInfo
}

func (s *Service) requestUrl(ctx context.Context, siteUrl string) (*common.SiteInfo, error) {
	now := time.Now()
	req, err := http.NewRequestWithContext(
		ctx, "GET", siteUrl, nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	var code int
	if resp != nil {
		code = resp.StatusCode
	} else {
		code = 500
	}
	return &common.SiteInfo{
		Url:     siteUrl,
		Status:  &code,
		Latency: time.Since(now),
		Error:   err,
	}, nil
}

func (s *Service) GetInfo(c *gin.Context) {
	s.requestStats.IncrementStats("getInfo")
	url := c.Query("url")
	if url == "" {
		c.IndentedJSON(http.StatusBadRequest, "Bad request: specify url")
		return
	}
	siteInfo, err := s.stats.GetSiteInfo(url)
	if err == nil {
		c.IndentedJSON(http.StatusOK, *siteInfo)
	} else {
		c.IndentedJSON(http.StatusNotFound, err.Error())
	}
}

func (s *Service) GetMin(c *gin.Context) {
	s.requestStats.IncrementStats("getMin")
	siteInfo, err := s.stats.GetMin()
	if err == nil {
		c.IndentedJSON(http.StatusOK, *siteInfo)
	} else {
		c.IndentedJSON(http.StatusNotFound, err.Error())
	}
}

func (s *Service) GetMax(c *gin.Context) {
	s.requestStats.IncrementStats("getMax")
	siteInfo, err := s.stats.GetMax()
	if err == nil {
		c.IndentedJSON(http.StatusOK, *siteInfo)
	} else {
		c.IndentedJSON(http.StatusNotFound, err.Error())
	}
}

func (s *Service) GetRequestStats(c *gin.Context) {
	data := s.requestStats.GetStats()
	c.IndentedJSON(http.StatusOK, data)
}
