// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateDaprComponentName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "-",
			Valid: false,
		},
		{
			Input: "9",
			Valid: false,
		},
		{
			Input: "a-",
			Valid: false,
		},
		{
			Input: "a--a",
			Valid: false,
		},
		{
			Input: "Cannothavecapitals",
			Valid: false,
		},
		{
			Input: "a",
			Valid: true,
		},
		{
			Input: "valid",
			Valid: true,
		},
		{
			Input: "valid12345678901234567890123456789012345678901234butverylong",
			Valid: true,
		},
		{
			Input: "invalid12345678901234567890123456789012345678901234567toolong",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := DaprComponentName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s", tc.Valid, valid, tc.Input)
		}
	}
}

func TestValidateSecretNames(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "-",
			Valid: false,
		},
		{
			Input: "9",
			Valid: true,
		},
		{
			Input: "a-",
			Valid: false,
		},
		{
			Input: "a.",
			Valid: false,
		},
		{
			Input: "a--a",
			Valid: true,
		},
		{
			Input: "Cannothavecapitals",
			Valid: false,
		},
		{
			Input: "no.dots.allowed",
			Valid: false,
		},
		{
			Input: "a",
			Valid: true,
		},
		{
			Input: "valid",
			Valid: true,
		},
		{
			Input: "valid12345678901234567890123456789012345678901234butverylong",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := SecretName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s: %+v", tc.Valid, valid, tc.Input, errors)
		}
	}
}

func TestValidateManagedEnvironmentStorageName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
		},
		{
			Input: "-",
		},
		{
			Input: "9",
		},
		{
			Input: "a-",
		},
		{
			Input: "a.",
		},
		{
			Input: "a--a",
		},
		{
			Input: "Cannothavecapitals",
		},
		{
			Input: "a-a",
			Valid: true,
		},
		{
			Input: "valid",
			Valid: true,
		},
		{
			Input: "invalid12345678901234567890123456",
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagedEnvironmentStorageName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s: %+v", tc.Valid, valid, tc.Input, errors)
		}
	}
}

func TestValidateManagedEnvironmentName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
		},
		{
			Input: "-",
		},
		{
			Input: "9",
		},
		{
			Input: "a-",
		},
		{
			Input: "a.",
		},
		{
			Input: "TooLonghijabcdefghijabcdefghijabcdefghijabcdefghijabcdefghija",
		},
		{
			Input: "a--a",
			Valid: true,
		},
		{
			Input: "Canhavecapitals",
			Valid: true,
		},
		{
			Input: "a-a",
			Valid: true,
		},
		{
			Input: "valid",
			Valid: true,
		},
		{
			Input: "MaxLengthjabcdefghijabcdefghijabcdefghijabcdefghijabcdefghij",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagedEnvironmentName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s: %+v", tc.Valid, valid, tc.Input, errors)
		}
	}
}

func TestValidateInitTimeout(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "5",
			Valid: false,
		},
		{
			Input: "m",
			Valid: false,
		},
		{
			Input: "6d",
			Valid: false,
		},
		{
			Input: "10s",
			Valid: true,
		},
		{
			Input: "1h",
			Valid: true,
		},
		{
			Input: "1200s",
			Valid: true,
		},
		{
			Input: "134m",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := InitTimeout(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s", tc.Valid, valid, tc.Input)
		}
	}
}

func TestContainerAppScaleRuleConcurrentRequests(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "5",
			Valid: true,
		},
		{
			Input: "m",
			Valid: false,
		},
		{
			Input: "6d",
			Valid: false,
		},
		{
			Input: "10s",
			Valid: false,
		},
		{
			Input: "1h",
			Valid: false,
		},
		{
			Input: "1200s",
			Valid: false,
		},
		{
			Input: "-1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ContainerAppScaleRuleConcurrentRequests(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s", tc.Valid, valid, tc.Input)
		}
	}
}

func TestValidateContainerAppName(t *testing.T) {
	// From Portal: Value must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character and cannot have '--'. The length must not be more than 32 characters.
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
		},
		{
			Input: "-",
		},
		{
			Input: "9",
		},
		{
			Input: "a-",
		},
		{
			Input: "a.",
		},
		{
			Input: "a--a",
		},
		{
			Input: "Cannothavecapitals",
		},
		{
			Input: "a-a",
			Valid: true,
		},
		{
			Input: "val-1-id",
			Valid: true,
		},
		{
			Input: "invalid12345678901234567890123456",
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ContainerAppName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s: %+v", tc.Valid, valid, tc.Input, errors)
		}
	}
}

func TestValidateContainerAppContainerName(t *testing.T) {
	// From Portal: Value must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character. The length must not be more than 46 characters.
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
		},
		{
			Input: "-",
		},
		{
			Input: "a-",
		},
		{
			Input: "a.",
		},
		{
			Input: "Cannothavecapitals",
		},
		{
			Input: "invalid1234567890123456789012345678901234567890",
		},
		{
			Input: "9",
			Valid: true,
		},
		{
			Input: "a",
			Valid: true,
		},
		{
			Input: "a--a",
			Valid: true,
		},
		{
			Input: "a-a",
			Valid: true,
		},
		{
			Input: "val-1-id",
			Valid: true,
		},
		{
			Input: "validbutverylong123456789012345678901234567890",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ContainerAppContainerName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t for %s: %+v", tc.Valid, valid, tc.Input, errors)
		}
	}
}
