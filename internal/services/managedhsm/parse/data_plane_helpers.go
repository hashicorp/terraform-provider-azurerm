// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type ManagedHSMDataPlaneEndpoint struct {
	ManagedHSMName string
	DomainSuffix   string
}

func (e ManagedHSMDataPlaneEndpoint) BaseURI() string {
	return fmt.Sprintf("https://%s.%s/", e.ManagedHSMName, e.DomainSuffix)
}

func ManagedHSMEndpoint(input string, domainSuffix *string) (*ManagedHSMDataPlaneEndpoint, error) {
	// NOTE: this function can be removed in 4.0
	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	return parseDataPlaneEndpoint(uri, domainSuffix)
}

func parseDataPlaneEndpoint(input *url.URL, domainSuffix *string) (*ManagedHSMDataPlaneEndpoint, error) {
	if input.Scheme != "https" {
		return nil, fmt.Errorf("expected the scheme to be `https` but got %q", input.Scheme)
	}

	hostname := strings.ToLower(input.Host)
	if port := input.Port(); port != "" {
		if port != "443" {
			return nil, fmt.Errorf("expected port to be '443' but got %q", port)
		}
		hostname = strings.TrimSuffix(hostname, fmt.Sprintf(":%s", input.Port()))
	}
	hostnameComponents := strings.Split(hostname, ".")
	if len(hostnameComponents) <= 2 {
		return nil, fmt.Errorf("expected the hostname to be in the format `{name}.{managedhsm}.{domain}` but got %q - please check this is a Managed HSM ID", hostname)
	}

	managedHSMName := hostnameComponents[0]
	if hostnameComponents[1] != "managedhsm" {
		return nil, fmt.Errorf("expected the hostname to contain `.managedhsm.` but got %q - please check this is a Managed HSM ID", hostname)
	}

	parsedDomainSuffix := strings.Join(hostnameComponents[1:], ".")
	if domainSuffix != nil && parsedDomainSuffix != *domainSuffix {
		// if we know the domain suffix, let's check that's what we're expecting
		return nil, fmt.Errorf("expected the hostname to end be in the format `{name}.%s` but got %q", *domainSuffix, hostname)
	}

	return &ManagedHSMDataPlaneEndpoint{
		ManagedHSMName: managedHSMName,
		DomainSuffix:   parsedDomainSuffix,
	}, nil
}

type dataPlaneResource struct {
	itemName    string
	itemVersion *string
}

func parseDataPlaneResource(input *url.URL, expectedType string, requireVersion bool) (*dataPlaneResource, error) {
	expectedSegments := 2
	expectedFormatExample := "/%s/{name}"
	if requireVersion {
		expectedSegments = 3
		expectedFormatExample = "/%s/{name}/{version}"
	}

	// then we want to check the path, which should be in the format described above
	path := strings.Split(strings.TrimPrefix(input.Path, "/"), "/")
	if len(path) != expectedSegments {
		return nil, fmt.Errorf("expected the path to be in the format %q but got %q", expectedFormatExample, input.Path)
	}

	nestedItemType := path[0]
	if nestedItemType != expectedType {
		return nil, fmt.Errorf("expected the Nested Item Type to be %q but got %q", expectedType, nestedItemType)
	}
	output := dataPlaneResource{
		itemName:    path[1],
		itemVersion: nil,
	}
	if err := validateSegment(output.itemName); err != nil {
		return nil, fmt.Errorf("expected the path to be in the format %q but %+v", expectedFormatExample, err)
	}

	if requireVersion {
		itemVersion := path[2]
		if err := validateSegment(itemVersion); err != nil {
			return nil, fmt.Errorf("expected the path to be in the format %q but %+v", expectedFormatExample, err)
		}
		output.itemVersion = pointer.To(itemVersion)
	}

	return &output, nil
}

func validateSegment(input string) error {
	val := strings.TrimSpace(input)
	if val == "" {
		return fmt.Errorf("unexpected empty value")
	}
	if val != input {
		return fmt.Errorf("contained extra whitespace")
	}

	return nil
}
