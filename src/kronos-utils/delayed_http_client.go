package lib

import (
	"net/http"
	"sync"
	"time"
)

type DelayedHttpClient struct {
	httpClient *http.Client
	delay      int64
	timestamp  int64
	resumeAt   int64

	doMutex sync.Mutex
}

func NewDelayedHttpClient(delayMillis int64) *DelayedHttpClient {
	return &DelayedHttpClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		delay:     delayMillis,
		timestamp: time.Now().UnixMilli(),
	}
}

func (hc *DelayedHttpClient) Do(req *http.Request) (*http.Response, error) {
	for {
		hc.doMutex.Lock()
		if hc.timestamp+hc.delay <= time.Now().UnixMilli() {
			hc.timestamp = time.Now().UnixMilli()
			hc.doMutex.Unlock()
			resp, err := hc.httpClient.Do(req)
			return resp, err
		} else {
			hc.doMutex.Unlock()
		}
	}
}
