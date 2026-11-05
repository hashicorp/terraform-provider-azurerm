// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/agentregistrationinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAutomationAccountFlatten(d *pluginsdk.ResourceData, id *automationaccount.AutomationAccountId, model *automationaccount.AutomationAccount, registration *agentregistrationinformation.AgentRegistration) error {
	d.Set("name", id.AutomationAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccessEnabled = *props.PublicNetworkAccess
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			skuName := ""
			if sku := props.Sku; sku != nil {
				skuName = string(sku.Name)
			}
			d.Set("sku_name", skuName)

			localAuthEnabled := true
			if val := props.DisableLocalAuth; val != nil && *val {
				localAuthEnabled = false
			}
			d.Set("local_authentication_enabled", localAuthEnabled)

			if err := d.Set("encryption", flattenEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}
			d.Set("hybrid_service_url", props.AutomationHybridServiceURL)

			identityVal, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", identityVal); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			d.Set("private_endpoint_connection", flattenPrivateEndpointConnections(props.PrivateEndpointConnections))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if registration != nil {
		d.Set("dsc_server_endpoint", registration.Endpoint)
		if keys := registration.Keys; keys != nil {
			d.Set("dsc_primary_access_key", keys.Primary)
			d.Set("dsc_secondary_access_key", keys.Secondary)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenEncryption(encryption *automationaccount.EncryptionProperties) []any {
	if encryption == nil || encryption.KeySource == nil || *encryption.KeySource != automationaccount.EncryptionKeySourceTypeMicrosoftPointKeyvault {
		return []any{}
	}

	keyVaultKeyId := ""
	userAssignedIdentityId := ""

	if keyProp := encryption.KeyVaultProperties; keyProp != nil {
		keyId, err := keyvault.NewNestedItemID(pointer.From(keyProp.KeyvaultUri), keyvault.NestedItemTypeKey, pointer.From(keyProp.KeyName), pointer.From(keyProp.KeyVersion))
		if err == nil {
			keyVaultKeyId = keyId.ID()
		}
	}

	if encryption.Identity != nil && encryption.Identity.UserAssignedIdentity != nil {
		userAssignedIdentityIdRaw := (*encryption.Identity.UserAssignedIdentity).(string)
		if userAssignedIdentityIdRaw != "" {
			uaiId, err := commonids.ParseUserAssignedIdentityIDInsensitively(userAssignedIdentityIdRaw)
			if err == nil {
				userAssignedIdentityId = uaiId.ID()
			}
		}
	}
	return []any{
		map[string]any{
			"key_vault_key_id":          keyVaultKeyId,
			"user_assigned_identity_id": userAssignedIdentityId,
		},
	}
}

func flattenPrivateEndpointConnections(input *[]automationaccount.PrivateEndpointConnection) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)
	for _, item := range *input {
		output = append(output, map[string]any{
			"id":   pointer.From(item.Id),
			"name": pointer.From(item.Name),
		})
	}
	return output
}
