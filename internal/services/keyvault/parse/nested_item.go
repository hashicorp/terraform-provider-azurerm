// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NestedItemObjectType enumerates the type of the "NestedItemType" value (e.g."keys", "secrets" or "certificates").
type NestedItemObjectType string

const (
	// KeyVaultObjectType Keys...
	NestedItemTypeKey NestedItemObjectType = "keys"
	// KeyVaultObjectType Secrets...
	NestedItemTypeSecret NestedItemObjectType = "secrets"
	// KeyVaultObjectType Certificates...
	NestedItemTypeCertificate NestedItemObjectType = "certificates"
	// KeyVaultObjectType Storage Keys...
	NestedItemTypeStorageKey NestedItemObjectType = "storage"
)

// PossibleNestedItemObjectTypeValues returns a string slice of possible "NestedItemObjectType" values.
func PossibleNestedItemObjectTypeValues() []string {
	return []string{string(NestedItemTypeKey), string(NestedItemTypeSecret), string(NestedItemTypeCertificate), string(NestedItemTypeStorageKey)}
}

var _ resourceids.Id = NestedItemId{}

type NestedItemId struct {
	KeyVaultBaseUrl string
	NestedItemType  NestedItemObjectType
	Name            string
	Version         string
}

func NewNestedItemID(keyVaultBaseUrl string, nestedItemType NestedItemObjectType, name, version string) (*NestedItemId, error) {
	keyVaultUrl, err := url.Parse(keyVaultBaseUrl)
	if err != nil || keyVaultBaseUrl == "" {
		return nil, fmt.Errorf("parsing %q: %+v", keyVaultBaseUrl, err)
	}
	// (@jackofallops) - Log Analytics service adds the port number to the API returns, so we strip it here
	if hostParts := strings.Split(keyVaultUrl.Host, ":"); len(hostParts) > 1 {
		keyVaultUrl.Host = hostParts[0]
	}

	if strings.Contains(strings.ToLower(keyVaultBaseUrl), ".managedhsm.") {
		return nil, fmt.Errorf("internal-error: Managed HSM IDs are not supported as Key Vault Nested Items")
	}

	return &NestedItemId{
		KeyVaultBaseUrl: keyVaultUrl.String(),
		NestedItemType:  nestedItemType,
		Name:            name,
		Version:         version,
	}, nil
}

func (id NestedItemId) ID() string {
	// example: https://tharvey-keyvault.vault.azure.net/type/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	segments := []string{
		strings.TrimSuffix(id.KeyVaultBaseUrl, "/"),
		string(id.NestedItemType),
		id.Name,
	}
	if id.Version != "" {
		segments = append(segments, id.Version)
	}
	return strings.TrimSuffix(strings.Join(segments, "/"), "/")
}

func (id NestedItemId) String() string {
	components := []string{
		fmt.Sprintf("Base Url %q", id.KeyVaultBaseUrl),
		fmt.Sprintf("Nested Item Type %q", string(id.NestedItemType)),
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Version %q", id.Version),
	}
	return fmt.Sprintf("Key Vault Nested Item %s", strings.Join(components, " / "))
}

func (id NestedItemId) VersionlessID() string {
	// example: https://tharvey-keyvault.vault.azure.net/type/bird
	segments := []string{
		strings.TrimSuffix(id.KeyVaultBaseUrl, "/"),
		string(id.NestedItemType),
		id.Name,
	}
	return strings.TrimSuffix(strings.Join(segments, "/"), "/")
}

// ParseNestedItemID parses a Key Vault Nested Item ID (such as a Certificate, Key or Secret)
// containing a version into a NestedItemId object
func ParseNestedItemID(input string) (*NestedItemId, error) {
	item, err := parseNestedItemId(input)
	if err != nil {
		return nil, err
	}

	if item.Version == "" {
		return nil, fmt.Errorf("expected a key vault versioned ID but no version information was found in: %q", input)
	}

	return item, nil
}

// ParseOptionallyVersionedNestedItemID parses a Key Vault Nested Item ID (such as a Certificate, Key or Secret)
// optionally containing a version into a NestedItemId object
func ParseOptionallyVersionedNestedItemID(input string) (*NestedItemId, error) {
	return parseNestedItemId(input)
}

func parseNestedItemId(id string) (*NestedItemId, error) {
	if strings.Contains(strings.ToLower(id), ".managedhsm.") {
		return nil, fmt.Errorf("internal-error: Managed HSM IDs are not supported as Key Vault Nested Items")
	}
	// versioned example: https://tharvey-keyvault.vault.azure.net/type/bird/fdf067c93bbb4b22bff4d8b7a9a56217
	// versionless example: https://tharvey-keyvault.vault.azure.net/type/bird/
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse azure key vault child ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 2 && len(components) != 3 {
		return nil, fmt.Errorf("key vault nested item should contain 2 or 3 segments, found %d segment(s) in %q", len(components), id)
	}

	version := ""
	if len(components) == 3 {
		version = components[2]
	}

	nestedItemObjectTypes := PossibleNestedItemObjectTypeValues()

	if !utils.SliceContainsValue(nestedItemObjectTypes, components[0]) {
		return nil, fmt.Errorf("key vault 'NestedItemType' should be one of: %s, got %q", strings.Join(nestedItemObjectTypes, ", "), components[0])
	}

	nestedItemObjectType := NestedItemObjectType(components[0])

	childId := NestedItemId{
		KeyVaultBaseUrl: fmt.Sprintf("https://%s/", idURL.Host),
		NestedItemType:  nestedItemObjectType,
		Name:            components[1],
		Version:         version,
	}

	return &childId, nil
}
