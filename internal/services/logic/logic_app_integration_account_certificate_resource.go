package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
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
			_, err := parse.IntegrationAccountCertificateID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountCertificateName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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

						"key_vault_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.VaultID,
						},

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

	id := parse.NewIntegrationAccountCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.CertificateName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_certificate", id.ID())
		}
	}

	parameters := logic.IntegrationAccountCertificate{
		IntegrationAccountCertificateProperties: &logic.IntegrationAccountCertificateProperties{},
	}

	if v, ok := d.GetOk("key_vault_key"); ok {
		parameters.IntegrationAccountCertificateProperties.Key = expandIntegrationAccountCertificateKeyVaultKey(v.([]interface{}))
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata, _ := pluginsdk.ExpandJsonFromString(v.(string))
		parameters.IntegrationAccountCertificateProperties.Metadata = metadata
	}

	if v, ok := d.GetOk("public_certificate"); ok {
		parameters.IntegrationAccountCertificateProperties.PublicCertificate = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IntegrationAccountName, id.CertificateName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountCertificateRead(d, meta)
}

func resourceLogicAppIntegrationAccountCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountCertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.CertificateName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if props := resp.IntegrationAccountCertificateProperties; props != nil {
		if err := d.Set("key_vault_key", flattenIntegrationAccountCertificateKeyVaultKey(props.Key)); err != nil {
			return fmt.Errorf("setting `key_vault_key`: %+v", err)
		}

		if props.Metadata != nil {
			metadataValue := props.Metadata.(map[string]interface{})
			metadataStr, _ := pluginsdk.FlattenJsonToString(metadataValue)
			d.Set("metadata", metadataStr)
		}

		d.Set("public_certificate", props.PublicCertificate)
	}

	return nil
}

func resourceLogicAppIntegrationAccountCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountCertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountCertificateID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationAccountName, id.CertificateName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountCertificateKeyVaultKey(input []interface{}) *logic.KeyVaultKeyReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := logic.KeyVaultKeyReference{
		KeyVault: &logic.KeyVaultKeyReferenceKeyVault{
			ID: utils.String(v["key_vault_id"].(string)),
		},
		KeyName: utils.String(v["key_name"].(string)),
	}

	if keyVersion := v["key_version"].(string); keyVersion != "" {
		result.KeyVersion = utils.String(keyVersion)
	}

	return &result
}

func flattenIntegrationAccountCertificateKeyVaultKey(input *logic.KeyVaultKeyReference) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var keyName string
	if input.KeyName != nil {
		keyName = *input.KeyName
	}

	var keyVaultId string
	if input.KeyVault != nil && input.KeyVault.ID != nil {
		keyVaultId = *input.KeyVault.ID
	}

	var keyVersion string
	if input.KeyVersion != nil {
		keyVersion = *input.KeyVersion
	}

	return []interface{}{
		map[string]interface{}{
			"key_name":     keyName,
			"key_vault_id": keyVaultId,
			"key_version":  keyVersion,
		},
	}
}
