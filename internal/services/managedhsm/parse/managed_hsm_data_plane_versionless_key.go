// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"
)

// NOTE: whilst we wouldn't normally prefix a struct and parser name with the package name, given the
// similarities between Key Vault and Managed HSM resources, we're intentionally doing that to make
// this clear, regardless of the import alias used for this package.

// ManagedHSMDataPlaneVersionlessKeyId defines the Data Plane ID for a Managed HSM Key without a Version.
// Example format: `https://{name}.{domainSuffix}/keys/{keyName}`
// Example value:  `https://example.managedhsm.azure.net/keys/bird`
type ManagedHSMDataPlaneVersionlessKeyId struct {
	// ManagedHSMName specifies the Name of this Managed HSM.
	ManagedHSMName string

	// DomainSuffix specifies the Domain Suffix used for Managed HSMs in the Azure Environment
	// where the Managed HSM exists - in the format `managedhsm.azure.net`.
	DomainSuffix string

	// KeyName specifies the name of this Managed HSM Key.
	KeyName string
}

// NewManagedHSMDataPlaneVersionlessKeyID returns a new instance of ManagedHSMDataPlaneVersionlessKeyId with the specified values.
func NewManagedHSMDataPlaneVersionlessKeyID(managedHsmName, domainSuffix, keyName string) ManagedHSMDataPlaneVersionlessKeyId {
	return ManagedHSMDataPlaneVersionlessKeyId{
		ManagedHSMName: managedHsmName,
		DomainSuffix:   domainSuffix,
		KeyName:        keyName,
	}
}

// ManagedHSMDataPlaneVersionlessKeyID parses the Data Plane Managed HSM Key ID without a Version.
func ManagedHSMDataPlaneVersionlessKeyID(input string, domainSuffix *string) (*ManagedHSMDataPlaneVersionlessKeyId, error) {
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}
	if domainSuffix != nil && !strings.HasPrefix(strings.ToLower(*domainSuffix), "managedhsm.") {
		return nil, fmt.Errorf("internal-error: the domainSuffix for Managed HSM %q didn't contain `managedhsm.`", *domainSuffix)
	}

	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	endpoint, err := parseDataPlaneEndpoint(uri, domainSuffix)
	if err != nil {
		// intentionally not wrapping this
		return nil, err
	}

	const requireVersion = false
	resource, err := parseDataPlaneResource(uri, "keys", requireVersion)
	if err != nil {
		// intentionally not wrapping this
		return nil, err
	}

	return &ManagedHSMDataPlaneVersionlessKeyId{
		ManagedHSMName: endpoint.ManagedHSMName,
		DomainSuffix:   endpoint.DomainSuffix,
		KeyName:        resource.itemName,
	}, nil
}

// BaseUri returns the Base URI for this Managed HSM Data Plane Versionless Key
func (id ManagedHSMDataPlaneVersionlessKeyId) BaseUri() string {
	return fmt.Sprintf("https://%s.%s/", id.ManagedHSMName, id.DomainSuffix)
}

// ID returns the full Resource ID for this Managed HSM Data Plane Versionless Key
func (id ManagedHSMDataPlaneVersionlessKeyId) ID() string {
	return fmt.Sprintf("https://%s.%s/keys/%s", id.ManagedHSMName, id.DomainSuffix, id.KeyName)
}

// String returns a human-readable description of this Managed HSM Key ID
func (id ManagedHSMDataPlaneVersionlessKeyId) String() string {
	components := []string{
		fmt.Sprintf("Managed HSM Name: %q", id.ManagedHSMName),
		fmt.Sprintf("Domain Suffix: %q", id.DomainSuffix),
		fmt.Sprintf("Key Name: %q", id.KeyName),
	}
	return fmt.Sprintf("Managed HSM Versionless Key (%s)", strings.Join(components, "\n"))
}
