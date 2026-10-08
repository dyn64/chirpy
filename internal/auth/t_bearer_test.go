package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBearerToken(t *testing.T) {
	cases := []struct {
		input      string
		expected   string
		want_match bool
	}{
		{
			input:      "test_token1234",
			expected:   "test_token1234",
			want_match: true,
		},
		{
			input:      "tezt",
			expected:   "tedt",
			want_match: false,
		},
	}

	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		bearer := fmt.Sprintf("Bearer %s", c.input)
		req.Header.Set("Authorization", bearer)

		token, err := GetBearerToken(req.Header)
		if err != nil {
			t.Fatalf("Error extracting token - %v", err)
		}
		diff := cmp.Diff(c.expected, token)

		if diff != "" {
			if c.want_match {
				t.Fatalf("%v", diff)
			}
		}
	}
}
