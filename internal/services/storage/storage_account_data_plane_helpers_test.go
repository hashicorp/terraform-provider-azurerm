// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"errors"
	"testing"
)

func TestConnectionError(t *testing.T) {
	testcases := []struct {
		Name        string
		Error       error
		ShouldMatch bool
	}{
		{
			Name:        "No Route TO Host",
			Error:       errors.New("dial tcp: connecting to example.blob.core.windows.net no route to host"),
			ShouldMatch: true,
		},
		{
			Name:        "DNS No Such Host",
			Error:       errors.New("dial tcp: lookup example.blob.core.windows.net on 10.0.0.1:53: no such host"),
			ShouldMatch: true,
		},
		{
			Name:        "Proxy Dropped",
			Error:       errors.New("EOF"),
			ShouldMatch: true,
		},
	}

	for _, tc := range testcases {
		if connectionError(tc.Error) != tc.ShouldMatch {
			t.Errorf("expected %s to match but it did not", tc.Name)
		}
	}
}
