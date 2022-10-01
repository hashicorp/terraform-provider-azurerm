package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type KeyRotationPolicyId struct {
	SubscriptionId     string
	ResourceGroup      string
	VaultName          string
	KeyName            string
	RotationpolicyName string
}

func NewKeyRotationPolicyID(subscriptionId, resourceGroup, vaultName, keyName string) KeyRotationPolicyId {
	return KeyRotationPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VaultName:      vaultName,
		KeyName:        keyName,
	}
}

func (id KeyRotationPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Key Name %q", id.KeyName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Key Rotation Policy", segmentsStr)
}

func (id KeyRotationPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/keys/%s/rotationpolicy"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.KeyName)
}

// KeyRotationPolicyID parses a KeyRotationPolicy ID into an KeyRotationPolicyId struct
func KeyRotationPolicyID(input string) (*KeyRotationPolicyId, error) {
	if !strings.HasSuffix(input, "/rotationpolicy") {
		return nil, fmt.Errorf("No rotation policy found in KeyRotationPolicy ID: %s", input)
	}

	id, err := resourceids.ParseAzureResourceID(strings.TrimSuffix(input, "/rotationpolicy"))
	if err != nil {
		return nil, err
	}

	resourceId := KeyRotationPolicyId{
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
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
