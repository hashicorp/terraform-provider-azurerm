// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_group_user

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/groupuser"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGroupUserCreate(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).GroupUsersClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := groupuser.NewGroupUserID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("group_name").(string), d.Get("user_id").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		exists, err := client.CheckEntityExists(ctx, id)
		if err != nil {
			if !response.WasNotFound(exists.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(exists.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_group_user", id.ID())
		}
	}

	if _, err := client.Create(ctx, id); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGroupUserRead(d, meta)
}
