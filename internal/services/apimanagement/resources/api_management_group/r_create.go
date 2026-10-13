// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/group"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GroupClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := group.NewGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	externalID := d.Get("external_id").(string)
	groupType := d.Get("type").(string)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_group", id.ID())
			}
		}
	}

	parameters := group.GroupCreateParameters{
		Properties: &group.GroupCreateParametersProperties{
			DisplayName: displayName,
			Description: pointer.To(description),
			ExternalId:  pointer.To(externalID),
			Type:        pointer.ToEnum[group.GroupType](groupType),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, group.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGroupRead(d, meta)
}
