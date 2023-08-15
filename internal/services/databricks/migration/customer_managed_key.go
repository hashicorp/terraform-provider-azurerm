// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = CustomerManagedKeyV0ToV1{}

type CustomerManagedKeyV0ToV1 struct {
}

func (c CustomerManagedKeyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"key_vault_key_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (c CustomerManagedKeyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// previously CustomerManagedKey used a virtual Resource ID - however this had a typo in the final segment `Mangaged`
		// whilst we could look to update this to remove the typo - since other Customer Managed Key resources use the
		// ID of the Parent resource being changed, we're going to update it to that instead.

		oldIdRaw := rawState["id"].(string)
		oldId, err := parseLegacyCustomerManagedKeyId(oldIdRaw)
		if err != nil {
			return nil, fmt.Errorf("parsing existing Resource ID %q: %+v", oldIdRaw, err)
		}

		newId := workspaces.NewWorkspaceID(oldId.SubscriptionId, oldId.ResourceGroup, oldId.CustomerMangagedKeyName)
		log.Printf("[DEBUG] Updating the Resource ID from %q to %q", oldIdRaw, newId.ID())
		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

type legacyCustomerManagedKeyId struct {
	SubscriptionId          string
	ResourceGroup           string
	CustomerMangagedKeyName string
}

func (id legacyCustomerManagedKeyId) String() string {
	segments := []string{
		fmt.Sprintf("Customer Mangaged Key Name %q", id.CustomerMangagedKeyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Customer Managed Key", segmentsStr)
}

func (id legacyCustomerManagedKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Databricks/customerMangagedKey/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CustomerMangagedKeyName)
}

func parseLegacyCustomerManagedKeyId(input string) (*legacyCustomerManagedKeyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := legacyCustomerManagedKeyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CustomerMangagedKeyName, err = id.PopSegment("customerMangagedKey"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
