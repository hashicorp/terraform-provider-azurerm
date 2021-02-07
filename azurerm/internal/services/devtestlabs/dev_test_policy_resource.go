package devtestlabs

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDevTestPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDevTestPolicyCreateUpdate,
		Read:   resourceArmDevTestPolicyRead,
		Update: resourceArmDevTestPolicyCreateUpdate,
		Delete: resourceArmDevTestPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(dtl.PolicyFactNameGalleryImage),
					string(dtl.PolicyFactNameLabPremiumVMCount),
					string(dtl.PolicyFactNameLabTargetCost),
					string(dtl.PolicyFactNameLabVMCount),
					string(dtl.PolicyFactNameLabVMSize),
					string(dtl.PolicyFactNameUserOwnedLabPremiumVMCount),
					string(dtl.PolicyFactNameUserOwnedLabVMCount),
					string(dtl.PolicyFactNameUserOwnedLabVMCountInSubnet),
				}, false),
			},

			"policy_set_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"lab_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"threshold": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"evaluator_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(dtl.AllowedValuesPolicy),
					string(dtl.MaxValuePolicy),
				}, false),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"fact_data": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDevTestPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Policy creation")

	name := d.Get("name").(string)
	policySetName := d.Get("policy_set_name").(string)
	labName := d.Get("lab_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, labName, policySetName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DevTest Policy %q (Policy Set %q / Lab %q / Resource Group %q): %s", name, policySetName, labName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dev_test_policy", *existing.ID)
		}
	}

	factData := d.Get("fact_data").(string)
	threshold := d.Get("threshold").(string)
	evaluatorType := d.Get("evaluator_type").(string)

	description := d.Get("description").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := dtl.Policy{
		Tags: tags.Expand(t),
		PolicyProperties: &dtl.PolicyProperties{
			FactName:      dtl.PolicyFactName(name),
			FactData:      utils.String(factData),
			Description:   utils.String(description),
			EvaluatorType: dtl.PolicyEvaluatorType(evaluatorType),
			Threshold:     utils.String(threshold),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, labName, policySetName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating DevTest Policy %q (Policy Set %q / Lab %q / Resource Group %q): %+v", name, policySetName, labName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, labName, policySetName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving DevTest Policy %q (Policy Set %q / Lab %q / Resource Group %q): %+v", name, policySetName, labName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read DevTest Policy %q (Policy Set %q / Lab %q / Resource Group %q) ID", name, policySetName, labName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDevTestPolicyRead(d, meta)
}

func resourceArmDevTestPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	labName := id.Path["labs"]
	policySetName := id.Path["policysets"]
	name := id.Path["policies"]

	read, err := client.Get(ctx, resourceGroup, labName, policySetName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] DevTest Policy %q was not found in Policy Set %q / Lab %q / Resource Group %q - removing from state!", name, policySetName, labName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DevTest Policy %q (Policy Set %q / Lab %q / Resource Group %q): %+v", name, policySetName, labName, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("policy_set_name", policySetName)
	d.Set("lab_name", labName)
	d.Set("resource_group_name", resourceGroup)

	if props := read.PolicyProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("fact_data", props.FactData)
		d.Set("evaluator_type", string(props.EvaluatorType))
		d.Set("threshold", props.Threshold)
	}

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceArmDevTestPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.PoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	labName := id.Path["labs"]
	policySetName := id.Path["policysets"]
	name := id.Path["policies"]

	read, err := client.Get(ctx, resourceGroup, labName, policySetName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] DevTest Policy %q was not found in Policy Set %q / Lab %q / Resource Group %q - assuming removed!", name, policySetName, labName, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving DevTest Policy %q (Policy Set %q / Lab %q / Resource Group %q): %+v", name, policySetName, labName, resourceGroup, err)
	}

	_, err = client.Delete(ctx, resourceGroup, labName, policySetName, name)
	if err != nil {
		return fmt.Errorf("Error deleting DevTest Policy %q (Policy Set %q / Lab %q / Resource Group %q): %+v", name, policySetName, labName, resourceGroup, err)
	}

	return err
}
