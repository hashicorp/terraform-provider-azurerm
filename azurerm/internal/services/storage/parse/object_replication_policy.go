package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ObjectReplicationPolicyId struct {
	SubscriptionId     string
	ResourceGroup      string
	StorageAccountName string
	Name               string
}

func NewObjectReplicationPolicyID(subscriptionId, resourceGroup, storageAccountName, name string) ObjectReplicationPolicyId {
	return ObjectReplicationPolicyId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		StorageAccountName: storageAccountName,
		Name:               name,
	}
}

func (id ObjectReplicationPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Object Replication Policy", segmentsStr)
}

func (id ObjectReplicationPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.Name)
}

// ObjectReplicationPolicyID parses a ObjectReplicationPolicy ID into an ObjectReplicationPolicyId struct
func ObjectReplicationPolicyID(input string) (*ObjectReplicationPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ObjectReplicationPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("objectReplicationPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
