// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/agentregistrationinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/automationaccount"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationAccountCreate,
		Read:   resourceAutomationAccountRead,
		Update: resourceAutomationAccountUpdate,
		Delete: resourceAutomationAccountDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := automationaccount.ParseAutomationAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(automationaccount.PossibleValuesForSkuNameEnum(), false),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"encryption": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},

						"key_source": {
							Type:       pluginsdk.TypeString,
							Optional:   true,
							Deprecated: "This field is now ignored and will be removed in the next major version of the Azure Provider, the `encryption` block can be omitted to disable encryption",
							ValidateFunc: validation.StringInSlice(
								[]string{
									string(automationaccount.EncryptionKeySourceTypeMicrosoftPointAutomation),
									string(automationaccount.EncryptionKeySourceTypeMicrosoftPointKeyvault),
								},
								false,
							),
						},

						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},
					},
				},
			},

			"local_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": commonschema.Tags(),

			"dsc_server_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"dsc_primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"dsc_secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"private_endpoint_connection": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"hybrid_service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutomationAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AutomationAccount
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := automationaccount.NewAutomationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_automation_account", id.ID())
	}

	identityVal, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	enc, err := expandEncryption(d.Get("encryption").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `encryption`: %v", err)
	}

	parameters := automationaccount.AutomationAccountCreateOrUpdateParameters{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &automationaccount.AutomationAccountCreateOrUpdateProperties{
			Encryption: enc,
			Sku: &automationaccount.Sku{
				Name: automationaccount.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: utils.Bool(d.Get("public_network_access_enabled").(bool)),
		},
	}

	parameters.Properties.DisableLocalAuth = utils.Bool(!d.Get("local_authentication_enabled").(bool))

	// for create account do not set identity property (even TypeNone is not allowed), or api will response error
	if identityVal.Type != identity.TypeNone {
		parameters.Identity = identityVal
	}
	if tagsVal := expandStringInterfaceMap(d.Get("tags").(map[string]interface{})); tagsVal != nil {
		parameters.Tags = &tagsVal
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAutomationAccountRead(d, meta)
}

func resourceAutomationAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AutomationAccount
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	enc, err := expandEncryption(d.Get("encryption").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `encryption`: %v", err)
	}

	parameters := automationaccount.AutomationAccountUpdateParameters{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: identity,
		Properties: &automationaccount.AutomationAccountUpdateProperties{
			Sku: &automationaccount.Sku{
				Name: automationaccount.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: utils.Bool(d.Get("public_network_access_enabled").(bool)),
			Encryption:          enc,
		},
	}

	if d.HasChange("local_authentication_enabled") {
		parameters.Properties.DisableLocalAuth = utils.Bool(!d.Get("local_authentication_enabled").(bool))
	}

	if tagsVal := tags.Expand(d.Get("tags").(map[string]interface{})); tagsVal != nil {
		parameters.Tags = tagsVal
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationAccountRead(d, meta)
}

func resourceAutomationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AutomationAccount
	registrationClient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	infoId := agentregistrationinformation.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)
	keysResp, err := registrationClient.Get(ctx, infoId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Agent Registration Info for %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Registration Info for %s: %+v", *id, err)
	}

	d.Set("name", id.AutomationAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
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

			identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			d.Set("private_endpoint_connection", flattenPrivateEndpointConnections(props.PrivateEndpointConnections))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	if model := keysResp.Model; model != nil {
		d.Set("dsc_server_endpoint", model.Endpoint)
		if keys := model.Keys; keys != nil {
			d.Set("dsc_primary_access_key", keys.Primary)
			d.Set("dsc_secondary_access_key", keys.Secondary)
		}
	}

	return nil
}

func resourceAutomationAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AutomationAccount
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandEncryption(input []interface{}) (*automationaccount.EncryptionProperties, error) {
	if len(input) == 0 {
		return &automationaccount.EncryptionProperties{
			KeySource: pointer.To(automationaccount.EncryptionKeySourceTypeMicrosoftPointAutomation),
		}, nil
	}

	v := input[0].(map[string]interface{})

	var id interface{}
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
		keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyIdStr)
		if err != nil {
			return nil, err
		}
		prop.KeyVaultProperties = &automationaccount.KeyVaultProperties{
			KeyName:     utils.String(keyId.Name),
			KeyVersion:  utils.String(keyId.Version),
			KeyvaultUri: utils.String(keyId.KeyVaultBaseUrl),
		}
	}
	return prop, nil
}

func flattenEncryption(encryption *automationaccount.EncryptionProperties) []interface{} {
	if encryption == nil || encryption.KeySource == nil || *encryption.KeySource != automationaccount.EncryptionKeySourceTypeMicrosoftPointKeyvault {
		return []interface{}{}
	}

	keyVaultKeyId := ""
	userAssignedIdentityId := ""

	if keyProp := encryption.KeyVaultProperties; keyProp != nil {
		keyId, err := keyVaultParse.NewNestedItemID(*keyProp.KeyvaultUri, keyVaultParse.NestedItemTypeKey, *keyProp.KeyName, *keyProp.KeyVersion)
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
	return []interface{}{
		map[string]interface{}{
			"key_vault_key_id":          keyVaultKeyId,
			"user_assigned_identity_id": userAssignedIdentityId,

			// TODO: remove this field in 4.x
			"key_source": "",
		},
	}
}

func flattenPrivateEndpointConnections(input *[]automationaccount.PrivateEndpointConnection) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, item := range *input {
		output = append(output, map[string]interface{}{
			"id":   pointer.From(item.Id),
			"name": pointer.From(item.Name),
		})
	}
	return output
}
