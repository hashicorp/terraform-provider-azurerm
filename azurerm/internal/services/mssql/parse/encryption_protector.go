package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EncryptionProtectorId struct {
	SubscriptionId          string
	ResourceGroup           string
	ServerName              string
	EncryptionProtectorName string
}

func NewEncryptionProtectorID(subscriptionId, resourceGroup, serverName, encryptionProtectorName string) EncryptionProtectorId {
	return EncryptionProtectorId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ServerName:              serverName,
		EncryptionProtectorName: encryptionProtectorName,
	}
}

func (id EncryptionProtectorId) String() string {
	segments := []string{
		fmt.Sprintf("Encryption Protector Name %q", id.EncryptionProtectorName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Encryption Protector", segmentsStr)
}

func (id EncryptionProtectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/encryptionProtector/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.EncryptionProtectorName)
}

// EncryptionProtectorID parses a EncryptionProtector ID into an EncryptionProtectorId struct
func EncryptionProtectorID(input string) (*EncryptionProtectorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EncryptionProtectorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.EncryptionProtectorName, err = id.PopSegment("encryptionProtector"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
