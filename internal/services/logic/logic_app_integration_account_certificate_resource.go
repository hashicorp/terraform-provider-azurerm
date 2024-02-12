// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppIntegrationAccountCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountCertificateCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountCertificateRead,
		Update: resourceLogicAppIntegrationAccountCertificateCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountCertificateDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationaccountcertificates.ParseCertificateID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountCertificateName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"integration_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"key_vault_key": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemName,
						},

						"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),

						"key_version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				AtLeastOneOf: []string{"public_certificate"},
			},

			"metadata": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"public_certificate": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"key_vault_key"},
			},
		},
	}
}

func resourceLogicAppIntegrationAccountCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountCertificateClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := integrationaccountcertificates.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_certificate", id.ID())
		}
	}

	parameters := integrationaccountcertificates.IntegrationAccountCertificate{
		Properties: integrationaccountcertificates.IntegrationAccountCertificateProperties{},
	}

	if v, ok := d.GetOk("key_vault_key"); ok {
		parameters.Properties.Key = expandIntegrationAccountCertificateKeyVaultKey(v.([]interface{}))
	}

	if v, ok := d.GetOk("metadata"); ok {
		parameters.Properties.Metadata = &v
	}

	if v, ok := d.GetOk("public_certificate"); ok {
		parameters.Properties.PublicCertificate = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountCertificateRead(d, meta)
}

func resourceLogicAppIntegrationAccountCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountcertificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if model := resp.Model; model != nil {
		props := model.Properties
		if err := d.Set("key_vault_key", flattenIntegrationAccountCertificateKeyVaultKey(props.Key)); err != nil {
			return fmt.Errorf("setting `key_vault_key`: %+v", err)
		}

		if props.Metadata != nil {
			d.Set("metadata", props.Metadata)
		}

		d.Set("public_certificate", props.PublicCertificate)

	}

	return nil
}

func resourceLogicAppIntegrationAccountCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountcertificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountCertificateKeyVaultKey(input []interface{}) *integrationaccountcertificates.KeyVaultKeyReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := integrationaccountcertificates.KeyVaultKeyReference{
		KeyVault: integrationaccountcertificates.KeyVaultKeyReferenceKeyVault{
			Id: utils.String(v["key_vault_id"].(string)),
		},
		KeyName: v["key_name"].(string),
	}

	if keyVersion := v["key_version"].(string); keyVersion != "" {
		result.KeyVersion = utils.String(keyVersion)
	}

	return &result
}

func flattenIntegrationAccountCertificateKeyVaultKey(input *integrationaccountcertificates.KeyVaultKeyReference) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var keyVaultId string
	if input.KeyVault.Id != nil {
		keyVaultId = *input.KeyVault.Id
	}

	var keyVersion string
	if input.KeyVersion != nil {
		keyVersion = *input.KeyVersion
	}

	return []interface{}{
		map[string]interface{}{
			"key_name":     input.KeyName,
			"key_vault_id": keyVaultId,
			"key_version":  keyVersion,
		},
	}
}
