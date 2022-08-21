package banisher

import (
	"testing"

	"github.com/renanqts/xdpfail2ban/pkg/config"
	"github.com/stretchr/testify/assert"
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
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected.bantime, actual.bantime)
			assert.Equal(t, tc.expected.findtime, actual.findtime)
		})
	}
}
