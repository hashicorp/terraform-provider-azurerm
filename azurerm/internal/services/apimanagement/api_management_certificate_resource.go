package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementCertificateCreateUpdate,
		Read:   resourceApiManagementCertificateRead,
		Update: resourceApiManagementCertificateCreateUpdate,
		Delete: resourceApiManagementCertificateDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	data := d.Get("data").(string)
	password := d.Get("password").(string)
	keyVaultSecretId := d.Get("key_vault_secret_id").(string)
	keyVaultIdentity := d.Get("key_vault_identity_client_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Certificate %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_certificate", *existing.ID)
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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("retrieving Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("cannot read ID for Certificate %q (Resource Group %q / API Management Service %q)", name, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

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
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	name := id.Name

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Certificate %q (Resource Group %q / API Management Service %q) was not found - removing from state!", name, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if props := resp.CertificateContractProperties; props != nil {
		if expiration := props.ExpirationDate; expiration != nil {
			formatted := expiration.Format(time.RFC3339)
			d.Set("expiration", formatted)
		}
		d.Set("subject", props.Thumbprint)
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
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	name := id.Name

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Certificate %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
		}
	}

	return nil
}
