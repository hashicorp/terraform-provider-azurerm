// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"reflect"
	"testing"
)

func TestManagedIdentityHeadersFromEnvironment(t *testing.T) {
	t.Setenv(managedIdentityHeaderEnvironmentVariable, "")
	if headers := ManagedIdentityHeadersFromEnvironment(); headers != nil {
		t.Fatalf("expected nil headers when `%s` is not set, got %#v", managedIdentityHeaderEnvironmentVariable, headers)
	}

	t.Setenv(managedIdentityHeaderEnvironmentVariable, "token")
	expected := map[string][]string{
		managedIdentityHeaderName: []string{"token"},
	}

	if headers := ManagedIdentityHeadersFromEnvironment(); !reflect.DeepEqual(headers, expected) {
		t.Fatalf("expected %#v, got %#v", expected, headers)
	}
}
