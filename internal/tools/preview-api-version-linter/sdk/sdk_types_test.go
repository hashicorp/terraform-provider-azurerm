package sdk

import (
	"regexp"
	"testing"
)

func TestSdkTypeServiceAndVersionRegex(t *testing.T) {
	testCases := []struct {
		regex           *regexp.Regexp
		input           string
		expectedService string
		expectedVersion string
	}{
		{
			regex:           GO_AZURE_SDK.ServiceAndVersionRegex,
			input:           "vendor/github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/tokens",
			expectedService: "containerregistry",
			expectedVersion: "2023-11-01-preview",
		},
		{
			regex:           KERMIT_SDK.ServiceAndVersionRegex,
			input:           "vendor/github.com/jackofallops/kermit/sdk/synapse/2020-08-01-preview",
			expectedService: "synapse",
			expectedVersion: "2020-08-01-preview",
		},
		{
			regex:           GIOVANNI_SDK.ServiceAndVersionRegex,
			input:           "vendor/github.com/jackofallops/giovanni/storage/2023-11-03-preview/file/shares",
			expectedService: "storage",
			expectedVersion: "2023-11-03-preview",
		},
		{
			regex:           AZURE_SDK_FOR_GO_TRACK_1.ServiceAndVersionRegex,
			input:           "vendor/github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview",
			expectedService: "resources",
			expectedVersion: "2021-06-01-preview",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			matches := tc.regex.FindStringSubmatch(tc.input)
			if len(matches) != 3 {
				t.Fatalf("expected 3 matches, got %d", len(matches))
			}
			service := matches[1]
			version := matches[2]

			if service != tc.expectedService {
				t.Errorf("expected service %s, got %s", tc.expectedService, service)
			}
			if version != tc.expectedVersion {
				t.Errorf("expected version %s, got %s", tc.expectedVersion, version)
			}
		})
	}
}
