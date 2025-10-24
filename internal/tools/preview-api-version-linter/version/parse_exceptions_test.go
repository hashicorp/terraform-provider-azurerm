package version

import (
	"regexp"
	"testing"
)

func TestParseException(t *testing.T) {
	testCases := []struct {
		name             string
		yamlBytes        []byte
		isHistorical     bool
		expectedVersions []Version
		expectedError    *regexp.Regexp
	}{
		{
			name: "Valid historical exceptions",
			yamlBytes: []byte(`
- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: blueprints
  version: 2018-11-01-preview

- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: containerregistry
  version: 2023-11-01-preview
`),
			isHistorical: true,
			expectedVersions: []Version{
				{Module: "github.com/hashicorp/go-azure-sdk/resource-manager", Service: "blueprints", Version: "2018-11-01-preview"},
				{Module: "github.com/hashicorp/go-azure-sdk/resource-manager", Service: "containerregistry", Version: "2023-11-01-preview"},
			},
		},
		{
			name: "Entries have to be sorted alphabetically by module, service and version",
			yamlBytes: []byte(`
- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: containerregistry
  version: 2023-11-01-preview

- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: blueprints
  version: 2018-11-01-preview
`),
			isHistorical:  true,
			expectedError: regexp.MustCompile("entries has to be sorted alphabetically by module, service and version"),
		},
		{
			name: "Unsupported modules",
			yamlBytes: []byte(`
- module: github.com/foo/bar
  service: baz
  version: 2020-01-01-preview
`),
			isHistorical:  true,
			expectedError: regexp.MustCompile("unsupported sdk module"),
		},
		{
			name: "Valid exceptions",
			yamlBytes: []byte(`
- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: blueprints
  version: 2018-11-01-preview
  stableVersionTargetDate: 2027-01-01
  responsibleIndividual: johndoe@microsoft.com

- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: containerregistry
  version: 2023-11-01-preview
  stableVersionTargetDate: 2027-01-01
  responsibleIndividual: github.com/gerrytan
`),
			isHistorical: false,
			expectedVersions: []Version{
				{Module: "github.com/hashicorp/go-azure-sdk/resource-manager", Service: "blueprints", Version: "2018-11-01-preview"},
				{Module: "github.com/hashicorp/go-azure-sdk/resource-manager", Service: "containerregistry", Version: "2023-11-01-preview"},
			},
		},
		{
			name: "Invalid stableVersionTargetDate",
			yamlBytes: []byte(`
- module: github.com/hashicorp/go-azure-sdk/resource-manager
  service: containerregistry
  version: 2020-01-01-preview
  stableVersionTargetDate: 01-01-2027
  responsibleIndividual: johndoe@email.com
`),
			isHistorical:  false,
			expectedError: regexp.MustCompile("invalid stableVersionTargetDate"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualVersions, actualError := parseExceptions(tc.yamlBytes, tc.isHistorical)

			if tc.expectedError != nil && actualError == nil {
				t.Fatalf("expected error: %q, got none", tc.expectedError)
			}
			if tc.expectedError == nil && actualError != nil {
				t.Fatalf("expected no error, got: %q", actualError)
			}
			if tc.expectedError != nil && actualError != nil && !tc.expectedError.MatchString(actualError.Error()) {
				t.Fatalf("expected error:\n\n%q\n\n, got:\n\n%q", tc.expectedError, actualError)
			}
			if len(tc.expectedVersions) != len(actualVersions) {
				t.Fatalf("expected %d versions, got: %d", len(tc.expectedVersions), len(actualVersions))
			}
			for i, expectedVersion := range tc.expectedVersions {
				if expectedVersion != actualVersions[i] {
					t.Errorf("expected version %d to be %+v, got: %+v", i, expectedVersion, actualVersions[i])
				}
			}
		})
	}
}

func TestValidResponsibleIndividual(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{input: "gerry.tan@microsoft.com", expected: true},
		{input: "github.com/gerrytan", expected: true},
		{input: "Gerry Tan", expected: false},
		{input: "gerrytan", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := validResponsibleIndividual(tc.input)
			if actual != tc.expected {
				t.Fatalf("expected %t, got: %t", tc.expected, actual)
			}
		})
	}
}
