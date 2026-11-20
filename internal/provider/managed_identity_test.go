// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"reflect"
	"testing"
)

func TestManagedIdentitySourceFromEnvironment(t *testing.T) {
	testCases := []struct {
		name     string
		env      map[string]string
		expected managedIdentitySource
	}{
		{
			name:     "no environment variables defaults to IMDS",
			env:      map[string]string{},
			expected: managedIdentitySourceDefaultToIMDS,
		},
		{
			name: "app service / container apps",
			env: map[string]string{
				identityEndpointEnvVar: "http://localhost/token",
				identityHeaderEnvVar:   "header-value",
			},
			expected: managedIdentitySourceAppService,
		},
		{
			name: "service fabric",
			env: map[string]string{
				identityEndpointEnvVar:         "http://localhost/token",
				identityHeaderEnvVar:           "header-value",
				identityServerThumbprintEnvVar: "thumbprint",
			},
			expected: managedIdentitySourceServiceFabric,
		},
		{
			name: "azure ml",
			env: map[string]string{
				msiEndpointEnvVar: "http://localhost/token",
				msiSecretEnvVar:   "secret",
			},
			expected: managedIdentitySourceAzureML,
		},
		{
			name: "cloud shell",
			env: map[string]string{
				msiEndpointEnvVar: "http://localhost/token",
			},
			expected: managedIdentitySourceCloudShell,
		},
		{
			name: "azure arc",
			env: map[string]string{
				identityEndpointEnvVar: "http://localhost/token",
				imdsEndpointEnvVar:     "http://localhost/imds",
			},
			expected: managedIdentitySourceAzureArc,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, envVar := range []string{identityEndpointEnvVar, identityHeaderEnvVar, identityServerThumbprintEnvVar, msiEndpointEnvVar, msiSecretEnvVar, imdsEndpointEnvVar} {
				t.Setenv(envVar, "")
			}
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			if actual := managedIdentitySourceFromEnvironment(); actual != testCase.expected {
				t.Fatalf("expected source %q but got %q", testCase.expected, actual)
			}
		})
	}
}

func TestManagedIdentityCustomHeaders(t *testing.T) {
	testCases := []struct {
		name     string
		env      map[string]string
		expected map[string][]string
	}{
		{
			name:     "no environment variables returns no headers",
			env:      map[string]string{},
			expected: nil,
		},
		{
			name: "app service / container apps",
			env: map[string]string{
				identityEndpointEnvVar: "http://localhost/token",
				identityHeaderEnvVar:   "header-value",
			},
			expected: map[string][]string{
				"X-IDENTITY-HEADER": {"header-value"},
			},
		},
		{
			name: "service fabric",
			env: map[string]string{
				identityEndpointEnvVar:         "http://localhost/token",
				identityHeaderEnvVar:           "header-value",
				identityServerThumbprintEnvVar: "thumbprint",
			},
			expected: map[string][]string{
				"Secret": {"header-value"},
			},
		},
		{
			name: "cloud shell returns no headers",
			env: map[string]string{
				msiEndpointEnvVar: "http://localhost/token",
			},
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, envVar := range []string{identityEndpointEnvVar, identityHeaderEnvVar, identityServerThumbprintEnvVar, msiEndpointEnvVar, msiSecretEnvVar, imdsEndpointEnvVar} {
				t.Setenv(envVar, "")
			}
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			if actual := ManagedIdentityCustomHeaders(); !reflect.DeepEqual(actual, testCase.expected) {
				t.Fatalf("expected headers %+v but got %+v", testCase.expected, actual)
			}
		})
	}
}
