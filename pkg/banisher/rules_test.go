package banisher

import (
	"testing"

	"github.com/renanqts/xdpfail2ban/pkg/config"
)

func TestTransformRules(t *testing.T) {
	tests := []struct {
		name     string
		rules    config.Rules
		expected TransformedRules
		err      error
	}{
		{
			name: "succeed",
			rules: config.Rules{
				Bantime:  "300s",
				Findtime: "120s",
			},
			expected: TransformedRules{
				bantime:  300000000000,
				findtime: 120000000000,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := transformRules(tc.rules)
			assertEqual(t, tc.err, err)
			assertEqual(t, tc.expected.bantime, actual.bantime)
			assertEqual(t, tc.expected.findtime, actual.findtime)
		})
	}
}
