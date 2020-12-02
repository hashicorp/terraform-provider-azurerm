package web

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCertificateCreateUpdate,
		Read:   resourceArmAppServiceCertificateRead,
		Update: resourceArmAppServiceCertificateCreateUpdate,
		Delete: resourceArmAppServiceCertificateDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CertificateID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"pfx_blob": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsBase64,
			},

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"key_vault_secret_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  azure.ValidateKeyVaultChildId,
				ConflictsWith: []string{"pfx_blob", "password"},
			},

			"hosting_environment_profile_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subject_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"host_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"issue_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmAppServiceCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	vaultClient := meta.(*clients.Client).KeyVault.VaultsClient
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Certificate creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	pfxBlob := d.Get("pfx_blob").(string)
	password := d.Get("password").(string)
	keyVaultSecretId := d.Get("key_vault_secret_id").(string)
	hostingEnvironmentProfileId := d.Get("hosting_environment_profile_id").(string)
	t := d.Get("tags").(map[string]interface{})

	if pfxBlob == "" && keyVaultSecretId == "" {
		return fmt.Errorf("Either `pfx_blob` or `key_vault_secret_id` must be set")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing App Service Certificate %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_certificate", *existing.ID)
		}
	}

	certificate := web.Certificate{
		CertificateProperties: &web.CertificateProperties{
			Password: utils.String(password),
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if len(hostingEnvironmentProfileId) > 0 {
		certificate.CertificateProperties.HostingEnvironmentProfile = &web.HostingEnvironmentProfile{
			ID: &hostingEnvironmentProfileId,
		}
	}

	if pfxBlob != "" {
		decodedPfxBlob, err := base64.StdEncoding.DecodeString(pfxBlob)
		if err != nil {
			return fmt.Errorf("Could not decode PFX blob: %+v", err)
		}
		certificate.CertificateProperties.PfxBlob = &decodedPfxBlob
	}

	if keyVaultSecretId != "" {
		parsedSecretId, err := azure.ParseKeyVaultChildID(keyVaultSecretId)
		if err != nil {
			return err
		}

		keyVaultBaseUrl := parsedSecretId.KeyVaultBaseUrl

		keyVaultId, err := azure.GetKeyVaultIDFromBaseUrl(ctx, vaultClient, keyVaultBaseUrl)
		if err != nil {
			return fmt.Errorf("Error retrieving the Resource ID for the Key Vault at URL %q: %s", keyVaultBaseUrl, err)
		}
		if keyVaultId == nil {
			return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", keyVaultBaseUrl)
		}

		certificate.CertificateProperties.KeyVaultID = keyVaultId
		certificate.CertificateProperties.KeyVaultSecretName = utils.String(parsedSecretId.Name)
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, certificate); err != nil {
		return fmt.Errorf("Error creating/updating App Service Certificate %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Certificate %q (Resource Group %q): %s", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service Certificate %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceCertificateRead(d, meta)
}

func resourceArmAppServiceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Certificate %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Certificate %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.CertificateProperties; props != nil {
		d.Set("friendly_name", props.FriendlyName)
		d.Set("subject_name", props.SubjectName)
		d.Set("host_names", props.HostNames)
		d.Set("issuer", props.Issuer)
		d.Set("issue_date", props.IssueDate.Format(time.RFC3339))
		d.Set("expiration_date", props.ExpirationDate.Format(time.RFC3339))
		d.Set("thumbprint", props.Thumbprint)

		if hep := props.HostingEnvironmentProfile; hep != nil {
			d.Set("hosting_environment_profile_id", hep.ID)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmAppServiceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Certificate %q (Resource Group %q)", id.Name, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting App Service Certificate %q (Resource Group %q): %s)", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
