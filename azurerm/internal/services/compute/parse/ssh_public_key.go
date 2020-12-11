package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SSHPublicKeyId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewSSHPublicKeyId(subscriptionId, resourceGroup, Name string) SSHPublicKeyId {
	return SSHPublicKeyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           Name,
	}
}

func (id SSHPublicKeyId) String() string {
	segments := []string{
		fmt.Sprintf("Key Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "SSH Public Key", segmentsStr)
}

func (id SSHPublicKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/sshPublicKeys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// SSHPublicKeyID parses a SSHPublicKey ID into an SSHPublicKeyId struct
func SSHPublicKeyID(input string) (*SSHPublicKeyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SSHPublicKeyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("sshPublicKeys"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
