package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/healthcareapis/mgmt/2018-08-20-preview/healthcareapis"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFhirApiService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFhirApiServiceCreateUpdate,
		Read:   resourceArmFhirApiServiceRead,
		Update: resourceArmFhirApiServiceCreateUpdate,
		Delete: resourceArmFhirApiServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "fhir",
			},

			"cosmodb_throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},

			"access_policy_object_ids": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmFhirApiServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).fhirApiServiceClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM FHIR API Service creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	kind := d.Get("kind").(string)
	cdba := int32(d.Get("cosmodb_throughput").(int))
	accessPolicyObjectIds := d.Get("access_policy_object_ids").([]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing FHIR API Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_fhir_api_service", *existing.ID)
		}
	}

	// create the cosmodb config object
	cosmoDbConfig := healthcareapis.ServiceCosmosDbConfigurationInfo{
		OfferThroughput: &cdba,
	}

	var svcAccessPolicyArray = []healthcareapis.ServiceAccessPolicyEntry{}
	for _, objectId := range accessPolicyObjectIds {
		objectIdMap := objectId.(map[string]interface{})
		objectIdsStr := objectIdMap["object_id"].(string)
		svcAccessPolicyObjectId := healthcareapis.ServiceAccessPolicyEntry{ObjectID: &objectIdsStr}
		svcAccessPolicyArray = append(svcAccessPolicyArray, svcAccessPolicyObjectId)
	}
	// create the service access policy array for the Service Properties
	properties := healthcareapis.ServicesProperties{
		AccessPolicies:        &svcAccessPolicyArray,
		CosmosDbConfiguration: &cosmoDbConfig,
	}

	fhirApiServiceDescription := healthcareapis.ServicesDescription{
		Location:   utils.String(location),
		Tags:       expandedTags,
		Kind:       &kind,
		Properties: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, fhirApiServiceDescription)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating FHIR API Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating FHIR API Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Retrieving FHIR API Service %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read FHIR API Service %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmFhirApiServiceRead(d, meta)
}

func resourceArmFhirApiServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).fhirApiServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	// ToDo Is services right here? Needs to get the account/instance name of fhir service
	name := id.Path["services"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] FHIR API Service %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure FHIR API Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmFhirApiServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).fhirApiServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error Parsing Azure Resource ID: %+v", err)
	}
	resGroup := id.ResourceGroup
	// ToDo Is services right here?
	name := id.Path["services"]
	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting FHIR API Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deleting FHIR API Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
