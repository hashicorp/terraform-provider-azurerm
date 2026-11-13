// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_service

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareServiceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := resource.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.ServicesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_healthcare_service", id.ID())
			}
		}
	}

	cosmosDbConfiguration, err := expandsCosmosDBConfiguration(d)
	if err != nil {
		return fmt.Errorf("expanding cosmosdb_configuration: %+v", err)
	}

	publicNetworkAccess := resource.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = resource.PublicNetworkAccessDisabled
	}

	payload := resource.ServicesDescription{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Kind:     resource.Kind(d.Get("kind").(string)),
		Properties: &resource.ServicesProperties{
			AccessPolicies:              expandAccessPolicyEntries(d),
			CosmosDbConfiguration:       cosmosDbConfiguration,
			CorsConfiguration:           expandCorsConfiguration(d),
			AuthenticationConfiguration: expandAuthentication(d),
			PublicNetworkAccess:         pointer.To(publicNetworkAccess),
		},
	}

	storageAcc, ok := d.GetOk("configuration_export_storage_account_name")
	if ok {
		payload.Properties.ExportConfiguration = &resource.ServiceExportConfigurationInfo{
			StorageAccountName: pointer.To(storageAcc.(string)),
		}
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	payload.Identity = expandedIdentity

	if d.IsNewResource() {
		if err := client.ServicesCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.ServicesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceHealthcareServiceRead(d, meta)
}
