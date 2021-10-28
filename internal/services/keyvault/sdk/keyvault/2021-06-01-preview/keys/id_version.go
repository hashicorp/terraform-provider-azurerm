package keys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VersionId struct {
	SubscriptionId string
	ResourceGroup  string
	VaultName      string
	KeyName        string
	Name           string
}

func NewVersionID(subscriptionId, resourceGroup, vaultName, keyName, name string) VersionId {
	return VersionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		KeyName:        keyName,
		Name:           name,
	}
}

func (id VersionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Version", segmentsStr)
}

func (id VersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/keys/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.KeyName, id.Name)
}

// ParseVersionID parses a Version ID into an VersionId struct
func ParseVersionID(input string) (*VersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VersionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
		return nil, err
	}
	if resourceId.KeyName, err = id.PopSegment("keys"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseVersionIDInsensitively parses an Version ID into an VersionId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseVersionID method should be used instead for validation etc.
func ParseVersionIDInsensitively(input string) (*VersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VersionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'vaults' segment
	vaultsKey := "vaults"
	for key := range id.Path {
		if strings.EqualFold(key, vaultsKey) {
			vaultsKey = key
			break
		}
	}
	if resourceId.VaultName, err = id.PopSegment(vaultsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'keys' segment
	keysKey := "keys"
	for key := range id.Path {
		if strings.EqualFold(key, keysKey) {
			keysKey = key
			break
		}
	}
	if resourceId.KeyName, err = id.PopSegment(keysKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'versions' segment
	versionsKey := "versions"
	for key := range id.Path {
		if strings.EqualFold(key, versionsKey) {
			versionsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(versionsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
