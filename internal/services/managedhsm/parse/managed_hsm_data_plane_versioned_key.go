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

// ManagedHSMDataPlaneVersionedKeyId defines the Data Plane ID for a Managed HSM Key with a Version.
// Example format: `https://{name}.{domainSuffix}/keys/{keyName}/{keyVersion}`
// Example value:  `https://example.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217`
type ManagedHSMDataPlaneVersionedKeyId struct {
	// ManagedHSMName specifies the Name of this Managed HSM.
	ManagedHSMName string

	// DomainSuffix specifies the Domain Suffix used for Managed HSMs in the Azure Environment
	// where the Managed HSM exists - in the format `managedhsm.azure.net`.
	DomainSuffix string

	// KeyName specifies the name of this Managed HSM Key.
	KeyName string

	// KeyVersion specifies the version of this Managed HSM Key.
	KeyVersion string
}

// NewManagedHSMDataPlaneVersionedKeyID returns a new instance of ManagedHSMDataPlaneVersionedKeyId with the specified values.
func NewManagedHSMDataPlaneVersionedKeyID(managedHsmName, domainSuffix, keyName, keyVersion string) ManagedHSMDataPlaneVersionedKeyId {
	return ManagedHSMDataPlaneVersionedKeyId{
		ManagedHSMName: managedHsmName,
		DomainSuffix:   domainSuffix,
		KeyName:        keyName,
		KeyVersion:     keyVersion,
	}
}

// ManagedHSMDataPlaneVersionedKeyID parses the Data Plane Managed HSM Key ID which requires a Version.
func ManagedHSMDataPlaneVersionedKeyID(input string, domainSuffix *string) (*ManagedHSMDataPlaneVersionedKeyId, error) {
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

	const requireVersion = true
	resource, err := parseDataPlaneResource(uri, "keys", requireVersion)
	if err != nil {
		// intentionally not wrapping this
		return nil, err
	}

	return &ManagedHSMDataPlaneVersionedKeyId{
		ManagedHSMName: endpoint.ManagedHSMName,
		DomainSuffix:   endpoint.DomainSuffix,
		KeyName:        resource.itemName,
		KeyVersion:     *resource.itemVersion,
	}, nil
}

// BaseUri returns the Base URI for this Managed HSM Data Plane Versioned Key
func (id ManagedHSMDataPlaneVersionedKeyId) BaseUri() string {
	return fmt.Sprintf("https://%s.%s/", id.ManagedHSMName, id.DomainSuffix)
}

// ID returns the full Resource ID for this Managed HSM Data Plane Versioned Key
func (id ManagedHSMDataPlaneVersionedKeyId) ID() string {
	return fmt.Sprintf("https://%s.%s/keys/%s/%s", id.ManagedHSMName, id.DomainSuffix, id.KeyName, id.KeyVersion)
}

// String returns a human-readable description of this Managed HSM Key ID
func (id ManagedHSMDataPlaneVersionedKeyId) String() string {
	components := []string{
		fmt.Sprintf("Managed HSM Name: %q", id.ManagedHSMName),
		fmt.Sprintf("Domain Suffix: %q", id.DomainSuffix),
		fmt.Sprintf("Key Name: %q", id.KeyName),
		fmt.Sprintf("Key Version: %q", id.KeyVersion),
	}
	return fmt.Sprintf("Managed HSM Key (%s)", strings.Join(components, "\n"))
}
