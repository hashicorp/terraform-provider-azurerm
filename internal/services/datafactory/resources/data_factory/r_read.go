// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.Factories
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := factories.ParseFactoryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, factories.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FactoryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			customerManagedKeyId := ""
			customerManagedKeyIdentityId := ""
			if enc := props.Encryption; enc != nil {
				if enc.VaultBaseURL != "" && enc.KeyName != "" && enc.KeyVersion != nil {
					version := ""
					if enc.KeyVersion != nil && *enc.KeyVersion != "" {
						version = *enc.KeyVersion
					}
					keyId, err := keyvault.NewNestedItemID(enc.VaultBaseURL, keyvault.NestedItemTypeKey, enc.KeyName, version)
					if err != nil {
						return fmt.Errorf("parsing Nested Item ID: %+v", err)
					}
					customerManagedKeyId = keyId.ID()
				}

				if encIdentity := enc.Identity; encIdentity != nil && pointer.From(encIdentity.UserAssignedIdentity) != "" {
					parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(encIdentity.UserAssignedIdentity))
					if err != nil {
						return err
					}
					customerManagedKeyIdentityId = parsed.ID()
				}
			}
			d.Set("customer_managed_key_id", customerManagedKeyId)
			d.Set("customer_managed_key_identity_id", customerManagedKeyIdentityId)

			globalParameters, err := flattenDataFactoryGlobalParameters(props.GlobalParameters)
			if err != nil {
				return fmt.Errorf("flattening `global_parameter`: %+v", err)
			}
			if err := d.Set("global_parameter", globalParameters); err != nil {
				return fmt.Errorf("setting `global_parameter`: %+v", err)
			}

			if err := d.Set("github_configuration", flattenGitHubRepoConfiguration(props.RepoConfiguration)); err != nil {
				return fmt.Errorf("setting `github_configuration`: %+v", err)
			}

			if err := d.Set("vsts_configuration", flattenVSTSRepoConfiguration(props.RepoConfiguration)); err != nil {
				return fmt.Errorf("setting `vsts_configuration`: %+v", err)
			}

			// The API defaults this to `true` but won't return it unless it's configured, so default it on
			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccessEnabled = *props.PublicNetworkAccess == factories.PublicNetworkAccessEnabled
			}
			d.Set("public_network_enabled", publicNetworkAccessEnabled)

			purviewId := ""
			if purview := props.PurviewConfiguration; purview != nil && purview.PurviewResourceId != nil {
				purviewId = *purview.PurviewResourceId
			}
			d.Set("purview_id", purviewId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	managedVirtualNetworkName, err := helpers.GetManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, id.SubscriptionId, id.ResourceGroupName, id.FactoryName)
	if err != nil {
		return err
	}
	d.Set("managed_virtual_network_enabled", managedVirtualNetworkName != nil)

	return nil
}
