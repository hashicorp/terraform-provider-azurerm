package validate

import (
	"math/rand"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestMonitorDiagnosticSettingName(t *testing.T) {
	cases := []struct {
		Name   string
		Errors int
	}{
		{
			Name:   "somename",
			Errors: 0,
		},
		{
			Name:   "",
			Errors: 1,
		},
		{
			Name:   RandStringBytes(261),
			Errors: 1,
		},
		{
			Name:   "some<name",
			Errors: 1,
		},
		{
			Name:   "some>name",
			Errors: 1,
		},
		{
			Name:   "some*name",
			Errors: 1,
		},
		{
			Name:   "some%name",
			Errors: 1,
		},
		{
			Name:   "some&name",
			Errors: 1,
		},
		{
			Name:   "some:name",
			Errors: 1,
		},
		{
			Name:   "some\\name",
			Errors: 1,
		},
		{
			Name:   "some?name",
			Errors: 1,
		},
		{
			Name:   "some+name",
			Errors: 1,
		},
		{
			Name:   "some/name",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := MonitorDiagnosticSettingName(tc.Name, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected Name to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
