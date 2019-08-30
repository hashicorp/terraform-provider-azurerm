package azurerm

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceCertificateCreateUpdate,
		Read:   resourceArmAppServiceCertificateRead,
		Update: resourceArmAppServiceCertificateCreateUpdate,
		Delete: resourceArmAppServiceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"pfx_blob": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validate.Base64String(),
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

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_vault_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
				ConflictsWith:    []string{"pfx_blob", "password"},
			},

			"key_vault_secret_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"pfx_blob", "password"},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmAppServiceCertificateCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.CertificatesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for App Service Certificate creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	password := d.Get("password").(string)
	t := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Certificate %q (Resource Group %q): %s", name, resourceGroup, err)
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

	if v, ok := d.GetOk("pfx_blob"); ok {
		pfxBlob, err := base64.StdEncoding.DecodeString(v.(string))
		if err != nil {
			return fmt.Errorf("Could not decode PFX blob: %+v", err)
		}
		certificate.CertificateProperties.PfxBlob = &pfxBlob
	}

	if v, ok := d.GetOk("key_vault_id"); ok {
		certificate.CertificateProperties.KeyVaultID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("key_vault_secret_name"); ok {
		certificate.CertificateProperties.KeyVaultSecretName = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, certificate); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read App Service Certificate %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceCertificateRead(d, meta)
}

func resourceArmAppServiceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.CertificatesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["certificates"]

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Certificate %q (Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Certificate %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.CertificateProperties; props != nil {
		d.Set("friendly_name", props.FriendlyName)
		d.Set("subject_name", props.SubjectName)
		d.Set("host_names", props.HostNames)
		d.Set("issuer", props.Issuer)
		d.Set("issue_date", props.IssueDate)
		d.Set("expiration_date", props.ExpirationDate)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("key_vault_id", props.KeyVaultID)
		d.Set("key_vault_secret_name", props.KeyVaultSecretName)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAppServiceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).web.CertificatesClient

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["certificates"]

	log.Printf("[DEBUG] Deleting App Service Certificate %q (Resource Group %q)", name, resourceGroup)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return err
		}
	}

	return nil
}
