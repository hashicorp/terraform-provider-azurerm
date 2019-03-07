package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBatchCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBatchCertificateCreate,
		Read:   resourceArmBatchCertificateRead,
		Update: resourceArmBatchCertificateUpdate,
		Delete: resourceArmBatchCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateAzureRMBatchCertificateName,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"certificate": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"format": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(batch.Cer),
					string(batch.Pfx),
				}, false),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true, // Required if `format` is "Cer"
				Sensitive: true,
			},
			"thumbprint": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"thumbprint_algorithm": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateAzureRMBatchCertificateThumbprint,
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceArmBatchCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchCertificateClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Batch certificate creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	name := d.Get("name").(string)
	certificate := d.Get("certificate").(string)
	format := d.Get("format").(string)
	password := d.Get("password").(string)
	thumbprint := d.Get("thumbprint").(string)
	thumbprintAlgorithm := d.Get("thumbprint_algorithm").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Batch Certificate %q (Resource Group %q, Account %q): %s", name, resourceGroupName, accountName, err)
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
		return fmt.Errorf("Error creating Batch certificate %q (Resource Group %q, Account %q): %+v", name, resourceGroupName, accountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Batch certificate %q (Resource Group %q, Account %q): %+v", name, resourceGroupName, accountName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch certificate %q (Resource Group %q), Account %q: %+v", name, resourceGroupName, accountName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch certificate %q (Resource Group %q, Account %q) ID", name, resourceGroupName, accountName)
	}

	d.SetId(*read.ID)

	return resourceArmBatchCertificateRead(d, meta)
}

func resourceArmBatchCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchCertificateClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	accountName := id.Path["batchAccounts"]
	name := id.Path["certificates"]
	resourceGroupName := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroupName, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Batch certificate %q was not found in Resource Group %q, Account %q - removing from state!", name, resourceGroupName, accountName)
			return nil
		}
		return fmt.Errorf("Error reading the state of Batch certificate %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("account_name", accountName)
	d.Set("resource_group_name", resourceGroupName)

	if format := resp.Format; format != "" {
		d.Set("format", format)
	}
	if publicData := resp.PublicData; publicData != nil {
		d.Set("public_data", publicData)
	}
	if thumbprint := resp.Thumbprint; thumbprint != nil {
		d.Set("thumbprint", thumbprint)
	}
	if thumbprintAlgorithm := resp.ThumbprintAlgorithm; thumbprintAlgorithm != nil {
		d.Set("thumbprint_algorithm", thumbprintAlgorithm)
	}

	return nil
}

func resourceArmBatchCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchCertificateClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Batch certificate update.")

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	accountName := id.Path["batchAccounts"]
	name := id.Path["certificates"]
	resourceGroupName := id.ResourceGroup

	certificate := d.Get("certificate").(string)
	format := d.Get("format").(string)
	password := d.Get("password").(string)
	thumbprint := d.Get("thumbprint").(string)
	thumbprintAlgorithm := d.Get("thumbprint_algorithm").(string)

	parameters := batch.CertificateCreateOrUpdateParameters{
		Name: &name,
		CertificateCreateOrUpdateProperties: &batch.CertificateCreateOrUpdateProperties{
			Data:                &certificate,
			Format:              batch.CertificateFormat(format),
			Password:            &password,
			Thumbprint:          &thumbprint,
			ThumbprintAlgorithm: &thumbprintAlgorithm,
		},
	}

	if _, err = client.Update(ctx, resourceGroupName, accountName, name, parameters, ""); err != nil {
		return fmt.Errorf("Error updating Batch certificate %q (Resource Group %q, Account %q): %+v", name, resourceGroupName, accountName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Batch certificate %q (Resource Group %q), Account %q: %+v", name, resourceGroupName, accountName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Batch certificate %q (Resource Group %q, Account %q) ID", name, resourceGroupName, accountName)
	}

	d.SetId(*read.ID)

	return resourceArmBatchCertificateRead(d, meta)
}

func resourceArmBatchCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchCertificateClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	accountName := id.Path["batchAccounts"]
	name := id.Path["certificates"]
	resourceGroupName := id.ResourceGroup

	future, err := client.Delete(ctx, resourceGroupName, accountName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Batch certificate %q (Resource Group %q, Account %q): %+v", name, resourceGroupName, accountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Batch certificate %q (Resource Group %q, Account %q): %+v", name, resourceGroupName, accountName, err)
		}
	}

	return nil
}

func validateAzureRMBatchCertificateName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w]+-[\w]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"must be made up of algorithm and thumbprint separated by a dash in %q: %q", k, value))
	}

	return warnings, errors
}
func validateAzureRMBatchCertificateThumbprint(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	switch value {
	case "SHA1":
	default:
		errors = append(errors, fmt.Errorf("currently required to be 'SHA1' in %q: %q", k, value))
	}
	return warnings, errors
}
