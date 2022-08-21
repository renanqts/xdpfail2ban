package banisher

import (
	"reflect"
	"regexp"
	"time"

	"github.com/renanqts/xdpfail2ban/pkg/config"
	"github.com/renanqts/xdpfail2ban/pkg/httpclient"
	"github.com/renanqts/xdpfail2ban/pkg/logger"
	"github.com/renanqts/xdpfail2ban/pkg/ttlmap"
)

// Banisher interface
type Banisher interface {
	Handler(string, string) error
}

// ViewedIP struct
type ViewedIP struct {
	viewed  time.Time
	counter int
	ban     bool
}

// ImpBanisher struct implements Banisher interface
type ImpBanisher struct {
	viewedIPs  *ttlmap.TTLMap
	rules      TransformedRules
	httpClient httpclient.Client
}

// New returns a Banisher
func New(c *config.Config) (Banisher, error) {
	rules, err := transformRules(c.Rules)
	if err != nil {
		return &ImpBanisher{}, err
	}

	b := &ImpBanisher{
		rules:      rules,
		httpClient: httpclient.New(c.XDPDropperURL),
	}

	b.viewedIPs = ttlmap.New(
		rules.bantime,
		func(key interface{}, value interface{}) {
			b.removeFromBan(key.(string), value.(ViewedIP))
		},
	)
	return b, nil
}

// Handler url requests
func (b *ImpBanisher) Handler(url, remoteIP string) error {
	logger.Debug.Printf("Banisher: New request %v", remoteIP)
	urlBytes := []byte(url)

	for _, reg := range b.rules.urlregexpBan {
		matched, err := regexp.Match(reg, urlBytes)
		if err != nil {
			return err
		}

		if matched {
			logger.Debug.Printf("Url ('%s') was matched by regexpBan: '%s'", url, reg)
			b.banController(remoteIP, func() time.Time { return time.Now() })
		}
	}

	return nil
}

func (b *ImpBanisher) banController(remoteIP string, now func() time.Time) {
	var viewedIP ViewedIP
	if value := b.viewedIPs.Get(remoteIP); value != nil {
		viewedIP = value.(ViewedIP)
	}

	if reflect.DeepEqual(viewedIP, ViewedIP{}) {
		logger.Debug.Printf("ip %s added into the viewed ips", remoteIP)
		b.viewedIPs.Add(remoteIP, ViewedIP{now(), 1, false})
	} else {
		if viewedIP.ban {
			// Check if IP is in the bantime,
			// in case yes, keep + count
			// in case no, delete it from the list
			if now().Before(viewedIP.viewed.Add(b.rules.bantime)) {
				b.viewedIPs.Add(remoteIP, ViewedIP{
					viewedIP.viewed,
					viewedIP.counter + 1,
					true,
				})
				b.drop("add", remoteIP)
				return
			}
			b.viewedIPs.Delete(remoteIP)
			logger.Debug.Println(remoteIP + " is no longer banned")
		} else if now().Before(viewedIP.viewed.Add(b.rules.findtime)) {
			// Check if IP counter >= maxretry
			// in case yes, block it
			// in case no, count
			if viewedIP.counter+1 >= b.rules.maxretry {
				b.viewedIPs.Add(remoteIP, ViewedIP{
					now(),
					viewedIP.counter + 1,
					true,
				})
				b.drop("add", remoteIP)
				logger.Debug.Println(remoteIP + " is now banned")
				return
			}
			b.viewedIPs.Add(remoteIP, ViewedIP{
				viewedIP.viewed,
				viewedIP.counter + 1,
				false,
			})
			logger.Debug.Printf("welcome back %s for the %d time", remoteIP, viewedIP.counter+1)
		} else {
			b.viewedIPs.Add(remoteIP, ViewedIP{
				now(),
				1,
				false,
			})
			logger.Debug.Printf("welcome back %s", remoteIP)
		}
	}
}

func (b *ImpBanisher) removeFromBan(ip string, ipViewed ViewedIP) {
	if ipViewed.ban {
		b.drop("remove", ip)
		logger.Debug.Printf("ip %s removed from ban. TTL expired", ip)
	}
}

func (b *ImpBanisher) drop(action, remoteIP string) {
	var err error
	body := struct{ IP string }{remoteIP}

	switch action {
	case "add":
		err = b.httpClient.Request(
			"/drop", "POST", 201, body,
		)
	case "remove":
		err = b.httpClient.Request(
			"/drop", "DELETE", 204, body,
		)
	}

	if err != nil {
		logger.Info.Printf(err.Error())
	}
}
