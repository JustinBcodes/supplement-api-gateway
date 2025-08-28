package proxy

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type TransportOpts struct {
	DialTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// SimpleReverseProxy forwards the request to a selected upstream URL string.
func SimpleReverseProxy(upstream string, rw http.ResponseWriter, req *http.Request) error {
	u, err := url.Parse(upstream)
	if err != nil {
		return err
	}
	// create new request to upstream
	upReq, err := http.NewRequestWithContext(req.Context(), req.Method, u.String()+req.URL.Path, req.Body)
	if err != nil {
		return err
	}
	upReq.Header = req.Header.Clone()
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(upReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	for k, vv := range resp.Header {
		for _, v := range vv {
			rw.Header().Add(k, v)
		}
	}
	rw.WriteHeader(resp.StatusCode)
	_, err = io.Copy(rw, resp.Body)
	return err
}


