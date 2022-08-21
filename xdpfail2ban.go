package xdpfail2ban

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/renanqts/xdpfail2ban/pkg/banisher"
	"github.com/renanqts/xdpfail2ban/pkg/config"
	"github.com/renanqts/xdpfail2ban/pkg/logger"
)

// XDPFail2Ban holds the necessary components of a Traefik plugin
type XDPFail2Ban struct {
	next     http.Handler
	name     string
	sync     sync.Mutex
	banisher banisher.Banisher
}

// CreateConfig populates the Config data object
func CreateConfig() *config.Config {
	return &config.Config{
		Rules: config.Rules{
			Bantime:  "300s",
			Findtime: "120s",
		},
	}
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *config.Config, name string) (http.Handler, error) {
	logger.SetLevel(config.LogLevel)

	banisher, err := banisher.New(config)
	if err != nil {
		return nil, err
	}

	logger.Info.Println("Plugin is up and running")
	return &XDPFail2Ban{
		next:     next,
		name:     name,
		banisher: banisher,
	}, nil
}

// Iterate over every headers to match the ones specified in the config and
// return nothing if regexp failed.
func (x *XDPFail2Ban) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger.Debug.Printf("New request: %v", req)

	remoteIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		logger.Debug.Println(remoteIP + " is not a valid IP or a IP/NET")
		return
	}
	url := req.URL.String()

	x.sync.Lock()
	defer x.sync.Unlock()

	err = x.banisher.Handler(url, remoteIP)
	if err != nil {
		logger.Info.Println("Request could not be handled " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
	x.next.ServeHTTP(rw, req)
}
