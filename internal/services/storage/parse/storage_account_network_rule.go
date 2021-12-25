package parse

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageAccountNetworkRuleId struct {
	StorageAccountId   *StorageAccountId
	IPRule             *storage.IPRule
	VirtualNetworkRule *storage.VirtualNetworkRule
	ResourceAccessRule *storage.ResourceAccessRule
}

func (id StorageAccountNetworkRuleId) ID() string {
	networkRuleId := ""
	if id.IPRule != nil && id.IPRule.IPAddressOrRange != nil {
		networkRuleId = fmt.Sprintf("ipAddressOrRange/%s", *id.IPRule.IPAddressOrRange)
	}

	if id.VirtualNetworkRule != nil && id.VirtualNetworkRule.VirtualNetworkResourceID != nil {
		subnetId := *id.VirtualNetworkRule.VirtualNetworkResourceID
		networkRuleId = fmt.Sprintf("subnetId/%s", strings.TrimPrefix(subnetId, "/"))
	}

	if id.ResourceAccessRule != nil && id.ResourceAccessRule.TenantID != nil && id.ResourceAccessRule.ResourceID != nil {
		resourceId := *id.ResourceAccessRule.ResourceID
		networkRuleId = fmt.Sprintf("tenantId/%s/resourceId/%s", *id.ResourceAccessRule.TenantID, strings.TrimPrefix(resourceId, "/"))
	}

	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s;%s"
	return fmt.Sprintf(fmtString, id.StorageAccountId.SubscriptionId, id.StorageAccountId.ResourceGroup, id.StorageAccountId.Name, networkRuleId)
}

func (id StorageAccountNetworkRuleId) String() string {
	segments := make([]string, 0)
	if id.IPRule != nil {
		segments = append(segments, fmt.Sprintf("IPAddressOrRange %q", *id.IPRule.IPAddressOrRange))
	}
	if id.VirtualNetworkRule != nil {
		segments = append(segments, fmt.Sprintf("SubnetID %q", *id.VirtualNetworkRule.VirtualNetworkResourceID))
	}
	if id.ResourceAccessRule != nil {
		segments = append(segments, fmt.Sprintf("TenantID %q", *id.ResourceAccessRule.TenantID))
		segments = append(segments, fmt.Sprintf("ResourceID %q", *id.ResourceAccessRule.ResourceID))
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s) for %s", "Network Rule", segmentsStr, id.StorageAccountId.String())
}

func StorageAccountNetworkRuleID(input string) (*StorageAccountNetworkRuleId, error) {
	segments := strings.Split(input, ";")
	if len(segments) != 2 {
		return nil, fmt.Errorf("storage network rule ID is composed as format `storageAccountId;networkRuleId`")
	}

	storageAccountId, err := StorageAccountID(segments[0])
	if err != nil {
		return nil, err
	}

	id := StorageAccountNetworkRuleId{
		StorageAccountId: storageAccountId,
	}

	if len(segments[1]) == 0 {
		return nil, fmt.Errorf("ID was missing the 'networkRuleId' element")
	}

	if strings.HasPrefix(segments[1], "ipAddressOrRange") {
		id.IPRule = &storage.IPRule{
			IPAddressOrRange: utils.String(strings.TrimPrefix(segments[1], "ipAddressOrRange/")),
		}
	}

	if strings.HasPrefix(segments[1], "subnetId") {
		id.VirtualNetworkRule = &storage.VirtualNetworkRule{
			VirtualNetworkResourceID: utils.String(strings.TrimPrefix(segments[1], "subnetId")),
		}
	}

	if strings.HasPrefix(segments[1], "tenantId") {
		var tenantId string
		var resourceId string

		components := strings.SplitN(segments[1], "/", 4)
		if len(components)%2 != 0 {
			return nil, fmt.Errorf("the number of path segments is not divisible by 2 in %q", segments[1])
		}

		for current := 0; current < len(components); current += 2 {
			key := components[current]
			value := components[current+1]
			if key == "" || value == "" {
				return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
			}

			switch {
			case key == "tenantId":
				tenantId = value
			case key == "resourceId":
				resourceId = "/" + value
			}
		}

		_, errors := validation.IsUUID(tenantId, "tenantId")
		if len(errors) != 0 {
			return nil, fmt.Errorf("%q should be a valid UUID: %q", tenantId, errors)
		}

		_, errors = azure.ValidateResourceID(resourceId, "resourceId")
		if len(errors) != 0 {
			return nil, fmt.Errorf("%q should be a valid Resource ID: %q", resourceId, errors)
		}

		id.ResourceAccessRule = &storage.ResourceAccessRule{
			TenantID:   utils.String(tenantId),
			ResourceID: utils.String(resourceId),
		}
	}

	return &id, nil
}
