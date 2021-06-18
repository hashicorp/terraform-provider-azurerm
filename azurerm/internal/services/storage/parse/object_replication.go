package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// This is manual for concat two ids are not supported in auto-generation

type ObjectReplicationId struct {
	SrcSubscriptionId     string
	SrcResourceGroup      string
	SrcStorageAccountName string
	SrcName               string
	DstSubscriptionId     string
	DstResourceGroup      string
	DstStorageAccountName string
	DstName               string
}

func NewObjectReplicationID(srcSubscriptionId, srcResourceGroup, strStorageAccountName, srcName, dstSubscriptionId, dstResourceGroup, dstStorageAccountName, dstName string) ObjectReplicationId {
	return ObjectReplicationId{
		SrcSubscriptionId:     srcSubscriptionId,
		SrcResourceGroup:      srcResourceGroup,
		SrcStorageAccountName: strStorageAccountName,
		SrcName:               srcName,
		DstSubscriptionId:     dstSubscriptionId,
		DstResourceGroup:      dstResourceGroup,
		DstStorageAccountName: dstStorageAccountName,
		DstName:               dstName,
	}
}

func (id ObjectReplicationId) String() string {
	segments := []string{
		fmt.Sprintf("Source Name %q", id.SrcName),
		fmt.Sprintf("Source Storage Account Name %q", id.SrcStorageAccountName),
		fmt.Sprintf("Source Resource Group %q", id.SrcResourceGroup),
		fmt.Sprintf("Source Subscription Id %q", id.SrcSubscriptionId),
		fmt.Sprintf("Destination Name %q", id.DstName),
		fmt.Sprintf("Destination Storage Account Name %q", id.DstStorageAccountName),
		fmt.Sprintf("Destination Resource Group %q", id.DstResourceGroup),
		fmt.Sprintf("Destination Subscription Id %q", id.DstSubscriptionId),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Object Replication", segmentsStr)
}

func (id ObjectReplicationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s;/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SrcSubscriptionId, id.SrcResourceGroup, id.SrcStorageAccountName, id.SrcName, id.DstSubscriptionId, id.DstResourceGroup, id.DstStorageAccountName, id.DstName)
}

func (id ObjectReplicationId) SourceObjectReplicationID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SrcSubscriptionId, id.SrcResourceGroup, id.SrcStorageAccountName, id.SrcName)
}

func (id ObjectReplicationId) DestinationObjectReplicationID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.DstSubscriptionId, id.DstResourceGroup, id.DstStorageAccountName, id.DstName)
}

// ObjectReplicationID parses a ObjectReplication ID into an ObjectReplicationId struct
func ObjectReplicationID(input string) (*ObjectReplicationId, error) {
	ids := strings.Split(input, ";")
	if len(ids) != 2 {
		return nil, fmt.Errorf("storage Object Replication Id is composed as format `sourceId;destinationId`")
	}
	srcId, err := azure.ParseAzureResourceID(ids[0])
	if err != nil {
		return nil, err
	}

	dstId, err := azure.ParseAzureResourceID(strings.TrimSuffix(ids[1], ";"))
	if err != nil {
		return nil, err
	}

	resourceId := ObjectReplicationId{
		SrcSubscriptionId: srcId.SubscriptionID,
		SrcResourceGroup:  srcId.ResourceGroup,
		DstSubscriptionId: dstId.SubscriptionID,
		DstResourceGroup:  dstId.ResourceGroup,
	}

	if resourceId.SrcSubscriptionId == "" || resourceId.DstSubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.SrcResourceGroup == "" || resourceId.DstResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SrcStorageAccountName, err = srcId.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.SrcName, err = srcId.PopSegment("objectReplicationPolicies"); err != nil {
		return nil, err
	}

	if err := srcId.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	if resourceId.DstStorageAccountName, err = dstId.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.DstName, err = dstId.PopSegment("objectReplicationPolicies"); err != nil {
		return nil, err
	}

	if err := dstId.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
