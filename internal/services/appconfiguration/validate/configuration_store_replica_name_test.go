// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestAppConfigurationReplicaName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "four",
			ErrCount: 0,
		},
		{
			Value:    "5five",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 1,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloWorld",
			ErrCount: 0,
		},
		{
			Value:    "helloworld12",
			ErrCount: 0,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd3324120",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqffewsqwcdw21ddwqwd332412020",
			ErrCount: 0,
		},
		{
			Value:    "qfvbdsbvipqdbwsbddbdcwqfjjfewsqwcdw21ddwqwd33241201",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := ConfigurationStoreReplicaName(tc.Value, "azurerm_app_configuration")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure App Configuration Replica Name to trigger a validation error: %v", tc)
		}
	}
}
