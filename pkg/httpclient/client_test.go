package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		expectedMethod     string
		statusCode         int
		expectedStatusCode int
		body               string
		expectedBody       string
		err                error
	}{
		{
			name:               "simple succeed request",
			method:             "POST",
			expectedMethod:     "POST",
			statusCode:         http.StatusCreated,
			expectedStatusCode: http.StatusCreated,
			body:               `{"foo":"bar"}`,
			expectedBody:       `{"foo":"bar"}`,
		},
		{
			name:               "unexpected status code",
			method:             "POST",
			expectedMethod:     "POST",
			statusCode:         http.StatusInternalServerError,
			expectedStatusCode: http.StatusCreated,
			body:               `{"foo":"bar"}`,
			expectedBody:       `{"foo":"bar"}`,
			err: fmt.Errorf(
				"xdpdropper request failed. Unexpected status code %d on %s. It should be %d",
				http.StatusInternalServerError,
				"POST",
				http.StatusCreated,
			),
		},
		{
			name:               "empty body",
			method:             "GET",
			expectedMethod:     "GET",
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "empty method",
			expectedMethod:     "GET",
			statusCode:         http.StatusOK,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "empty status code",
			method:         "GET",
			expectedMethod: "GET",
			err: fmt.Errorf(
				"xdpdropper request failed. Unexpected status code %d on %s. It should be %d",
				http.StatusOK,
				"GET",
				0,
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, tc.expectedMethod, r.Method)

					if len(tc.expectedBody) == 0 {
						reqBody, err := io.ReadAll(r.Body)
						assert.Nil(t, err)
						var body string
						err = json.Unmarshal(reqBody, &body)
						assert.Nil(t, err)
						assert.Equal(t, tc.expectedBody, body)
					}

					if tc.statusCode != 0 {
						w.WriteHeader(tc.statusCode)
					}
				}),
			)
			defer server.Close()
			c := New(server.URL)
			err := c.Request("/foobar", tc.method, tc.expectedStatusCode, tc.body)
			assert.Equal(t, tc.err, err)
		})
	}

}
