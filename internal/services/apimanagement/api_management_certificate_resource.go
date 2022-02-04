package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementCertificateCreateUpdate,
		Read:   resourceApiManagementCertificateRead,
		Update: resourceApiManagementCertificateCreateUpdate,
		Delete: resourceApiManagementCertificateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"data": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Sensitive:     true,
				ValidateFunc:  validation.StringIsBase64,
				AtLeastOneOf:  []string{"data", "key_vault_secret_id"},
				ConflictsWith: []string{"key_vault_secret_id", "key_vault_identity_client_id"},
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"data"},
			},

			"key_vault_secret_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  keyVaultValidate.NestedItemIdWithOptionalVersion,
				AtLeastOneOf:  []string{"data", "key_vault_secret_id"},
				ConflictsWith: []string{"data", "password"},
			},

			"key_vault_identity_client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"key_vault_secret_id"},
			},

			"expiration": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subject": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiManagementCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	data := d.Get("data").(string)
	password := d.Get("password").(string)
	keyVaultSecretId := d.Get("key_vault_secret_id").(string)
	keyVaultIdentity := d.Get("key_vault_identity_client_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_certificate", id.ID())
		}
	}

	parameters := apimanagement.CertificateCreateOrUpdateParameters{
		CertificateCreateOrUpdateProperties: &apimanagement.CertificateCreateOrUpdateProperties{},
	}

	if keyVaultSecretId != "" {
		parsedSecretId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultSecretId)
		if err != nil {
			return err
		}

		parameters.KeyVault = &apimanagement.KeyVaultContractCreateProperties{
			SecretIdentifier: utils.String(parsedSecretId.ID()),
		}

		if keyVaultIdentity != "" {
			parameters.KeyVault.IdentityClientID = utils.String(keyVaultIdentity)
		}
	}

	if data != "" {
		parameters.Data = utils.String(data)
		parameters.Password = utils.String(password)
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.Name, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementCertificateRead(d, meta)
}

func resourceApiManagementCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	if props := resp.CertificateContractProperties; props != nil {
		if expiration := props.ExpirationDate; expiration != nil {
			formatted := expiration.Format(time.RFC3339)
			d.Set("expiration", formatted)
		}
		d.Set("subject", props.Subject)
		d.Set("thumbprint", props.Thumbprint)

		if keyvault := props.KeyVault; keyvault != nil {
			d.Set("key_vault_secret_id", keyvault.SecretIdentifier)
			d.Set("key_vault_identity_client_id", keyvault.IdentityClientID)
		}
	}

	return nil
}

func resourceApiManagementCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
