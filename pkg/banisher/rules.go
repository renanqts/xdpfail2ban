package banisher

import (
	"time"

	"github.com/renanqts/xdpfail2ban/pkg/config"
	"github.com/renanqts/xdpfail2ban/pkg/logger"
)

// TransformedRules transformed config.Rules struct
type TransformedRules struct {
	bantime      time.Duration
	findtime     time.Duration
	urlregexpBan []string
	maxretry     int
}

// transformRule morph a Rules object into a TransformedRules
func transformRules(r config.Rules) (TransformedRules, error) {
	bantime, err := time.ParseDuration(r.Bantime)
	if err != nil {
		return TransformedRules{}, err
	}
	logger.Info.Printf("Bantime: %s", bantime)

	findtime, err := time.ParseDuration(r.Findtime)
	if err != nil {
		return TransformedRules{}, err
	}
	logger.Info.Printf("Findtime: %s", findtime)

	var regexpBan []string
	for _, rg := range r.Urlregexps {
		logger.Info.Printf("using rule %q", rg.Regexp)
		regexpBan = append(regexpBan, rg.Regexp)
	}

	rules := TransformedRules{
		bantime:      bantime,
		findtime:     findtime,
		urlregexpBan: regexpBan,
		maxretry:     r.Maxretry,
	}

	logger.Info.Printf("Loaded rules: '%+v'", rules)
	return rules, nil
}
