package batch

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2019-08-01/batch"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceBatchCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceBatchCertificateCreate,
		Read:   resourceBatchCertificateRead,
		Update: resourceBatchCertificateUpdate,
		Delete: resourceBatchCertificateDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CertificateID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureRMBatchAccountName,
			},

			// TODO: make this case sensitive once this API bug has been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/5574
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"certificate": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 10000),
			},

			"format": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(batch.Cer),
					string(batch.Pfx),
				}, false),
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true, // Required if `format` is "Pfx"
				Sensitive: true,
			},

			"thumbprint": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"thumbprint_algorithm": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringInSlice([]string{"SHA1"}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"public_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBatchCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch certificate creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	certificate := d.Get("certificate").(string)
	format := d.Get("format").(string)
	password := d.Get("password").(string)
	thumbprint := d.Get("thumbprint").(string)
	thumbprintAlgorithm := d.Get("thumbprint_algorithm").(string)
	name := thumbprintAlgorithm + "-" + thumbprint

	if err := validateBatchCertificateFormatAndPassword(format, password); err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Batch Certificate %q (Account %q / Resource Group %q): %s", name, accountName, resourceGroupName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_batch_certificate", *existing.ID)
		}
	}
	certificateProperties := batch.CertificateCreateOrUpdateProperties{
		Data:                &certificate,
		Format:              batch.CertificateFormat(format),
		Thumbprint:          &thumbprint,
		ThumbprintAlgorithm: &thumbprintAlgorithm,
	}
	if password != "" {
		certificateProperties.Password = &password
	}
	parameters := batch.CertificateCreateOrUpdateParameters{
		Name:                                &name,
		CertificateCreateOrUpdateProperties: &certificateProperties,
	}

	future, err := client.Create(ctx, resourceGroupName, accountName, name, parameters, "", "")
	if err != nil {
		return fmt.Errorf("Error creating Batch certificate %q (Account %q / Resource Group %q): %+v", name, accountName, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Batch certificate %q (Account %q / Resource Group %q): %+v", name, accountName, resourceGroupName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch certificate %q (Account %q / Resource Group %q): %+v", name, accountName, resourceGroupName, err)
	}
	d.SetId(*read.ID)
	return resourceBatchCertificateRead(d, meta)
}

func resourceBatchCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Batch certificate %q was not found in Account %q / Resource Group %q - removing from state!", id.Name, id.BatchAccountName, id.ResourceGroup)
			return nil
		}
		return fmt.Errorf("Error retrieving Batch Certificate %q (Account %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.CertificateProperties; props != nil {
		d.Set("format", props.Format)
		d.Set("public_data", props.PublicData)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("thumbprint_algorithm", props.ThumbprintAlgorithm)
	}

	return nil
}

func resourceBatchCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch certificate update.")

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	certificate := d.Get("certificate").(string)
	format := d.Get("format").(string)
	password := d.Get("password").(string)
	thumbprint := d.Get("thumbprint").(string)
	thumbprintAlgorithm := d.Get("thumbprint_algorithm").(string)

	if err := validateBatchCertificateFormatAndPassword(format, password); err != nil {
		return err
	}

	parameters := batch.CertificateCreateOrUpdateParameters{
		Name: &id.Name,
		CertificateCreateOrUpdateProperties: &batch.CertificateCreateOrUpdateProperties{
			Data:                &certificate,
			Format:              batch.CertificateFormat(format),
			Password:            &password,
			Thumbprint:          &thumbprint,
			ThumbprintAlgorithm: &thumbprintAlgorithm,
		},
	}

	if _, err = client.Update(ctx, id.ResourceGroup, id.BatchAccountName, id.Name, parameters, ""); err != nil {
		return fmt.Errorf("Error updating Batch certificate %q (Account %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch Certificate %q (Account %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Batch certificate %q (Account: %q, Resource Group %q) ID", id.Name, id.BatchAccountName, id.ResourceGroup)
	}

	return resourceBatchCertificateRead(d, meta)
}

func resourceBatchCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Batch Certificate %q (Account %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Batch Certificate %q (Account %q / Resource Group %q): %+v", id.Name, id.BatchAccountName, id.ResourceGroup, err)
		}
	}

	return nil
}

func validateBatchCertificateFormatAndPassword(format string, password string) error {
	if format == "Pfx" && password == "" {
		return fmt.Errorf("Batch Certificate Password is required when Format is `Pfx`")
	}
	if format == "Cer" && password != "" {
		return fmt.Errorf(" Batch Certificate Password must not be specified when Format is `Cer`")
	}
	return nil
}
