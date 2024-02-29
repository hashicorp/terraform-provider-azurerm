// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ManagedHSMKeyID{}

type ManagedHSMKeyID struct {
	HSMBaseUrl string
	Name       string
	Version    string
}

func NewManagedHSMKeyID(keyVaultBaseUrl string, name, version string) (*ManagedHSMKeyID, error) {
	keyVaultUrl, err := url.Parse(keyVaultBaseUrl)
	if err != nil || keyVaultBaseUrl == "" {
		return nil, fmt.Errorf("parsing %q: %+v", keyVaultBaseUrl, err)
	}
	// (@jackofallops) - Log Analytics service adds the port number to the API returns, so we strip it here
	if hostParts := strings.Split(keyVaultUrl.Host, ":"); len(hostParts) > 1 {
		keyVaultUrl.Host = hostParts[0]
	}

	if !strings.Contains(strings.ToLower(keyVaultBaseUrl), ".managedhsm.") {
		return nil, fmt.Errorf("internal-error: only support Managed HSM IDs as Nested Items")
	}

	return &ManagedHSMKeyID{
		HSMBaseUrl: keyVaultUrl.String(),
		Name:       name,
		Version:    version,
	}, nil
}

func (id ManagedHSMKeyID) ID() string {
	// example: https://tharvey-managedhsm.managedhsm.azure.net/keys/key-bird/fdf067c93bbb4b22bff4d8b7a9a56217
	segments := []string{
		strings.TrimSuffix(id.HSMBaseUrl, "/"),
		"keys",
		id.Name,
	}
	if id.Version != "" {
		segments = append(segments, id.Version)
	}
	return strings.TrimSuffix(strings.Join(segments, "/"), "/")
}

func (id ManagedHSMKeyID) String() string {
	components := []string{
		fmt.Sprintf("Base Url %q", id.HSMBaseUrl),
		fmt.Sprintf("Type %q", "keys"),
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Version %q", id.Version),
	}
	return fmt.Sprintf("Managed HSM Nested Item %s", strings.Join(components, " / "))
}

func (id ManagedHSMKeyID) VersionlessID() string {
	// example: https://tharvey-managedhsm.managedhsm.azure.net/keys/key-bird
	segments := []string{
		strings.TrimSuffix(id.HSMBaseUrl, "/"),
		"keys",
		id.Name,
	}
	return strings.TrimSuffix(strings.Join(segments, "/"), "/")
}

// ParseNestedItemID parses a managed HSM Nested Item ID (such as a Certificate, Key or Secret)
// containing a version into a NestedItemId object
func ParseNestedItemID(input string) (*ManagedHSMKeyID, error) {
	item, err := parseNestedItemId(input)
	if err != nil {
		return nil, err
	}

	if item.Version == "" {
		return nil, fmt.Errorf("expected a managed HSM versioned ID but no version information was found in: %q", input)
	}

	return item, nil
}

// ParseOptionallyVersionedNestedItemID parses a managed HSM Nested Item ID (such as a Certificate, Key or Secret)
// optionally containing a version into a NestedItemId object
func ParseOptionallyVersionedNestedItemID(input string) (*ManagedHSMKeyID, error) {
	return parseNestedItemId(input)
}

func parseNestedItemId(id string) (*ManagedHSMKeyID, error) {
	if !strings.Contains(strings.ToLower(id), ".managedhsm.") {
		return nil, fmt.Errorf("internal-error: only support Managed HSM IDs as managed HSM Nested Items")
	}
	// versioned example: https://tharvey-managedhsm.managedhsm.azure.net/keys/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	// versionless example: https://tharvey-managedhsm.managedhsm.azure.net/keys/bird/
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse azure managed HSM child ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 2 && len(components) != 3 {
		return nil, fmt.Errorf("managed HSM nested item should contain 2 or 3 segments, found %d segment(s) in %q", len(components), id)
	}

	if components[0] != "keys" {
		return nil, fmt.Errorf("expected a managed HSM nested item of type 'keys' but got %q", components[0])
	}

	version := ""
	if len(components) == 3 {
		version = components[2]
	}

	childId := ManagedHSMKeyID{
		HSMBaseUrl: fmt.Sprintf("https://%s/", idURL.Host),
		Name:       components[1],
		Version:    version,
	}

	return &childId, nil
}
