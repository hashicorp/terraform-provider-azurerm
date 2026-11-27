// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package database_migration_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2021-06-30/serviceresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	databasemigrationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabaseMigrationServiceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := databasemigrationClient.From(meta).ServicesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := serviceresource.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.ServicesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_database_migration_service", id.ID())
		}
	}

	parameters := serviceresource.DataMigrationService{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &serviceresource.DataMigrationServiceProperties{
			VirtualSubnetId: d.Get("subnet_id").(string),
		},
		Sku: &serviceresource.ServiceSku{
			Name: pointer.To(d.Get("sku_name").(string)),
		},
		Kind: pointer.To("Cloud"), // currently only "Cloud" is supported, hence hardcode here
	}
	if t, ok := d.GetOk("tags"); ok {
		parameters.Tags = tags.Expand(t.(map[string]any))
	}

	if err := client.ServicesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}
	return resourceDatabaseMigrationServiceRead(d, meta)
}
