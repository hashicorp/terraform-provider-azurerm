package validate

import (
	"testing"
)

func TestSecretRotationInterval(t *testing.T) {
	cases := []struct {
		SecretRotationInterval string
		Errors                 int
	}{
		{
			SecretRotationInterval: "",
			Errors:                 1,
		},
		{
			SecretRotationInterval: "2m",
			Errors:                 0,
		},
		{
			SecretRotationInterval: "1md",
			Errors:                 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.SecretRotationInterval, func(t *testing.T) {
			_, errors := SecretRotationInterval(tc.SecretRotationInterval, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected SecretRotationInterval to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
