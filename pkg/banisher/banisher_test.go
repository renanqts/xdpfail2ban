package banisher

import (
	"testing"
	"time"

	"github.com/renanqts/xdpfail2ban/pkg/ttlmap"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal()
}

type MockHTTPClient struct {
	FakeRequest func(string, string, int, interface{}) error
}

func (m *MockHTTPClient) Request(context, method string, statusCode int, body interface{}) error {
	return nil
}

func TestBanController(t *testing.T) {
	b := ImpBanisher{
		httpClient: &MockHTTPClient{
			FakeRequest: func(context string, method string, statusCode int, body interface{}) error {
				return nil
			},
		},
	}

	callback := func(interface{}, interface{}) {}

	tests := []struct {
		name             string
		shouldExist      bool
		rules            TransformedRules
		ip               string
		time             time.Time
		viewedIP         ViewedIP
		expectedViewedIP ViewedIP
	}{
		{
			name:        "increase counter",
			shouldExist: true,
			rules: TransformedRules{
				bantime:  300000000000,
				findtime: 120000000000,
				maxretry: 2,
			},
			ip: "1.1.1.1",
			time: time.Date(
				2016,
				time.April,
				18,
				23,
				15,
				8,
				4,
				time.UTC,
			),
			viewedIP: ViewedIP{
				viewed: time.Date(
					2016,
					time.April,
					18,
					23,
					15,
					8,
					4,
					time.UTC,
				),
			},
			expectedViewedIP: ViewedIP{
				viewed: time.Date(
					2016,
					time.April,
					18,
					23,
					15,
					8,
					4,
					time.UTC,
				),
				counter: 1,
				ban:     false,
			},
		},
		{
			name:        "must ban",
			shouldExist: true,
			rules: TransformedRules{
				bantime:  300000000000,
				findtime: 120000000000,
				maxretry: 2,
			},
			ip: "1.1.1.1",
			time: time.Date(
				2016,
				time.April,
				18,
				23,
				15,
				8,
				4,
				time.UTC,
			),
			viewedIP: ViewedIP{
				viewed: time.Date(
					2016,
					time.April,
					18,
					23,
					15,
					8,
					4,
					time.UTC,
				),
				counter: 1,
				ban:     false,
			},
			expectedViewedIP: ViewedIP{
				viewed: time.Date(
					2016,
					time.April,
					18,
					23,
					15,
					8,
					4,
					time.UTC,
				),
				counter: 2,
				ban:     true,
			},
		},
		{
			name:        "remove from ban",
			shouldExist: false,
			rules: TransformedRules{
				bantime:  300000000000,
				findtime: 120000000000,
				maxretry: 2,
			},
			ip: "1.1.1.1",
			time: time.Date(
				2016,
				time.April,
				18,
				20,
				15,
				8,
				4,
				time.UTC,
			),
			viewedIP: ViewedIP{
				viewed: time.Date(
					2016,
					time.April,
					18,
					4,
					15,
					8,
					4,
					time.UTC,
				),
				counter: 2,
				ban:     true,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b.viewedIPs = ttlmap.New(tc.rules.bantime, callback)
			b.rules = tc.rules
			if (ViewedIP{} != tc.viewedIP) {
				b.viewedIPs.Add(tc.ip, tc.viewedIP)
			}
			b.banController(tc.ip, func() time.Time { return tc.time })
			var actualIPView ViewedIP
			value := b.viewedIPs.Get(tc.ip)
			if tc.shouldExist {
				actualIPView = value.(ViewedIP)
				assertEqual(t, tc.expectedViewedIP.viewed, actualIPView.viewed)
				assertEqual(t, tc.expectedViewedIP.counter, actualIPView.counter)
				assertEqual(t, tc.expectedViewedIP.ban, actualIPView.ban)
			} else {
				if value != nil {
					t.Fatal()
				}
			}
		})
	}
}
