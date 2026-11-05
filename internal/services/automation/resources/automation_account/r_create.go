// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.AutomationAccount
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := automationaccount.NewAutomationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_account", id.ID())
		}
	}

	identityVal, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	enc, err := expandEncryption(d.Get("encryption").([]any))
	if err != nil {
		return fmt.Errorf("expanding `encryption`: %v", err)
	}

	parameters := automationaccount.AutomationAccountCreateOrUpdateParameters{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &automationaccount.AutomationAccountCreateOrUpdateProperties{
			Encryption: enc,
			Sku: &automationaccount.Sku{
				Name: automationaccount.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: pointer.To(d.Get("public_network_access_enabled").(bool)),
		},
	}

	parameters.Properties.DisableLocalAuth = pointer.To(!d.Get("local_authentication_enabled").(bool))

	// for create account do not set identity property (even TypeNone is not allowed), or api will response error
	if identityVal.Type != identity.TypeNone {
		parameters.Identity = identityVal
	}
	if tagsVal := helpers.ExpandStringInterfaceMap(d.Get("tags").(map[string]any)); tagsVal != nil {
		parameters.Tags = &tagsVal
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceAutomationAccountRead(d, meta)
}

func expandEncryption(input []any) (*automationaccount.EncryptionProperties, error) {
	if len(input) == 0 {
		return &automationaccount.EncryptionProperties{
			KeySource: pointer.To(automationaccount.EncryptionKeySourceTypeMicrosoftPointAutomation),
		}, nil
	}

	v := input[0].(map[string]any)

	var id any
	id, ok := v["user_assigned_identity_id"].(string)
	if !ok {
		return nil, fmt.Errorf("read encryption user identity id error")
	}
	prop := &automationaccount.EncryptionProperties{
		KeySource: pointer.To(automationaccount.EncryptionKeySourceTypeMicrosoftPointKeyvault),
	}
	if id != "" {
		prop.Identity = &automationaccount.EncryptionPropertiesIdentity{
			UserAssignedIdentity: &id,
		}
	}

	if keyIdStr := v["key_vault_key_id"].(string); keyIdStr != "" {
		keyId, err := keyvault.ParseNestedItemID(keyIdStr, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}
		prop.KeyVaultProperties = &automationaccount.KeyVaultProperties{
			KeyName:     pointer.To(keyId.Name),
			KeyVersion:  pointer.To(keyId.Version),
			KeyvaultUri: pointer.To(keyId.KeyVaultBaseURL),
		}
	}
	return prop, nil
}
