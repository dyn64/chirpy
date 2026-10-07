package auth

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestEncodeToken(t *testing.T) {
	id := uuid.New()
	cases := []struct {
		input    uuid.UUID
		expected uuid.UUID
	}{
		{
			input:    id,
			expected: id,
		},
	}

	for _, c := range cases {
		token, err := MakeJWT(id, "Allyourbase", time.Minute*20)
		if err != nil {
			t.Fatalf("Error making token - %v", err)
		}

		ret_token, err := ValidateJWT(token, "Allyourbase")
		if err != nil {
			t.Fatalf("Error validating token - %v", err)
		}
		diff := cmp.Diff(c.input, ret_token)
		if diff != "" {
			t.Fatalf("%v", diff)
		}
	}
}
