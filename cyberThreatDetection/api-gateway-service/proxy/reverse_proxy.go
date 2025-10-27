package proxy

import (
	"io"
	"net/http"
	"strings"
)

type ReverseProxy struct {
	logServiceURL		string
	threatServiceURL	string
}

func NewReverseProxy(logServiceURL, threatServiceURL string) *ReverseProxy {
	return &ReverseProxy{
		logServiceURL: logServiceURL,
		threatServiceURL: threatServiceURL,
	}
}

func (p *ReverseProxy) Route(w http.ResponseWriter, r *http.Request) {
	path :=r.URL.Path

	var targetURL string
	if strings.HasPrefix(path, "/api/logs"){
		targetURL = p.logServiceURL
	} else if strings.HasPrefix(path, "/api/threats") {
		targetURL = p.threatServiceURL
	} else {
		http.Error(w, `{"success":false,"error":{"code":404,"message":"Not found"}}`, http.StatusNotFound)
		return
	}

	p.ProxyRequest(w, r, targetURL)
}

func (p *ReverseProxy) ProxyRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	url := targetURL + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	ProxyReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, `{"success":false,"error":{"code":500,"message":"Failed to create request"}}`, http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			ProxyReq.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(ProxyReq)
	if err != nil {
		http.Error(w, `{"success":false,"error":{"code":502,"message":"Service unavailable"}}`, http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}