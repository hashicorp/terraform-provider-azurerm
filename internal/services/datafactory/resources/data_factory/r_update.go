// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedvirtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.Factories
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := factories.NewFactoryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, factories.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_factory", id.ID())
			}
		}
	}

	publicNetworkAccess := factories.PublicNetworkAccessEnabled
	enabled := d.Get("public_network_enabled").(bool)
	if !enabled {
		publicNetworkAccess = factories.PublicNetworkAccessDisabled
	}

	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	normalizedLocation := location.Normalize(d.Get("location").(string))
	payload := factories.Factory{
		Location: pointer.To(normalizedLocation),
		Properties: &factories.FactoryProperties{
			PublicNetworkAccess: &publicNetworkAccess,
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if purviewId, ok := d.GetOk("purview_id"); ok {
		payload.Properties.PurviewConfiguration = &factories.PurviewConfiguration{
			PurviewResourceId: pointer.To(purviewId.(string)),
		}
	}

	if keyVaultKeyID, ok := d.GetOk("customer_managed_key_id"); ok {
		keyVaultKey, err := keyvault.ParseNestedItemID(keyVaultKeyID.(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
		if err != nil {
			return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}

		payload.Properties.Encryption = &factories.EncryptionConfiguration{
			VaultBaseURL: keyVaultKey.KeyVaultBaseURL,
			KeyName:      keyVaultKey.Name,
			KeyVersion:   &keyVaultKey.Version,
			Identity:     helpers.ExpandDataFactoryEncryptionIdentity(d.Get("customer_managed_key_identity_id").(string)),
		}
	}

	globalParameters, err := expandDataFactoryGlobalParameters(d.Get("global_parameter").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	payload.Properties.GlobalParameters = globalParameters

	if _, err := client.CreateOrUpdate(ctx, id, payload, factories.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	githubConfiguration := expandGitHubRepoConfiguration(d.Get("github_configuration").([]any))
	if githubConfiguration != nil {
		repoUpdate := factories.FactoryRepoUpdate{
			FactoryResourceId: pointer.To(id.ID()),
			RepoConfiguration: githubConfiguration,
		}
		locationId := factories.NewLocationID(id.SubscriptionId, normalizedLocation)
		if _, err := client.ConfigureFactoryRepo(ctx, locationId, repoUpdate); err != nil {
			return fmt.Errorf("configuring Repository for %s: %+v", locationId, err)
		}
	}
	vstsConfiguration := expandVSTSRepoConfiguration(d.Get("vsts_configuration").([]any))
	if vstsConfiguration != nil {
		repoUpdate := factories.FactoryRepoUpdate{
			FactoryResourceId: pointer.To(id.ID()),
			RepoConfiguration: vstsConfiguration,
		}
		locationId := factories.NewLocationID(id.SubscriptionId, normalizedLocation)
		if _, err := client.ConfigureFactoryRepo(ctx, locationId, repoUpdate); err != nil {
			return fmt.Errorf("configuring Repository for %s: %+v", locationId, err)
		}
	}

	if d.Get("managed_virtual_network_enabled").(bool) {
		networkPayload := managedvirtualnetworks.ManagedVirtualNetworkResource{
			Properties: managedvirtualnetworks.ManagedVirtualNetwork{},
		}
		managedNetworkId := managedvirtualnetworks.NewManagedVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName, "default")
		if _, err := managedVirtualNetworksClient.CreateOrUpdate(ctx, managedNetworkId, networkPayload, managedvirtualnetworks.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("creating virtual network for %s: %+v", id, err)
		}
	}

	return resourceDataFactoryRead(d, meta)
}
