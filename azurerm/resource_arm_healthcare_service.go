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

func resourceArmHealthcareService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHealthcareServiceCreateUpdate,
		Read:   resourceArmHealthcareServiceRead,
		Update: resourceArmHealthcareServiceCreateUpdate,
		Delete: resourceArmHealthcareServiceDelete,

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

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "fhir",
			},

			"cosmosdb_throughput": {
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
			/*
				"cors_configuration": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"origins": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"headers": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"methods": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"max_age": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"allow_credentials": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
			*/
			"tags": tagsSchema(),
		},
	}
}

func resourceArmHealthcareServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).healthcare.HealthcareServiceClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Healthcare Service creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	kind := d.Get("kind").(string)
	cdba := int32(d.Get("cosmosdb_throughput").(int))
	accessPolicyObjectIds := d.Get("access_policy_object_ids").([]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Healthcare Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_healthcare_service", *existing.ID)
		}
	}

	var svcAccessPolicyArray []healthcareapis.ServiceAccessPolicyEntry
	for _, objectId := range accessPolicyObjectIds {
		objectIdMap := objectId.(map[string]interface{})
		objectIdsStr := objectIdMap["object_id"].(string)
		svcAccessPolicyObjectId := healthcareapis.ServiceAccessPolicyEntry{ObjectID: &objectIdsStr}
		svcAccessPolicyArray = append(svcAccessPolicyArray, svcAccessPolicyObjectId)
	}

	healthcareServiceDescription := healthcareapis.ServicesDescription{
		Location: utils.String(location),
		Tags:     expandedTags,
		Kind:     &kind,
		Properties: &healthcareapis.ServicesProperties{
			AccessPolicies: &svcAccessPolicyArray,
			CosmosDbConfiguration: &healthcareapis.ServiceCosmosDbConfigurationInfo{
				OfferThroughput: &cdba,
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, healthcareServiceDescription)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Retrieving Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Healthcare Service %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHealthcareServiceRead(d, meta)
}

func resourceArmHealthcareServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).healthcare.HealthcareServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["services"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Healthcare Service %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Healthcare Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if kind := resp.Kind; kind != nil {
		d.Set("kind", kind)
	}
	if properties := resp.Properties; properties != nil {
		if err := d.Set("access_policy_object_ids", azure.FlattenHealthcareAccessPolicies(properties.AccessPolicies)); err != nil {
			return fmt.Errorf("Error setting 'access_policy_object_ids': %+v", err)
		}
		if config := properties.CosmosDbConfiguration; config != nil {
			d.Set("cosmosdb_throughput", config.OfferThroughput)
		}
		if authConfig := properties.AuthenticationConfiguration; authConfig != nil {
			d.Set("authentication_configuration_authority", authConfig.Authority)
			d.Set("authentication_configuration_audience", authConfig.Audience)
			d.Set("authentication_configuration_smart_proxy_enabled", authConfig.SmartProxyEnabled)
		}
		corsOutput := make([]interface{}, 0)
		if corsConfig := properties.CorsConfiguration; corsConfig != nil {
			output := make(map[string]interface{})
			if corsConfig.Origins != nil {
				output["origins"] = *corsConfig.Origins
			}
			if corsConfig.Headers != nil {
				output["headers"] = *corsConfig.Headers
			}
			if corsConfig.Methods != nil {
				output["methods"] = *corsConfig.Methods
			}
			if corsConfig.MaxAge != nil {
				output["max_age"] = *corsConfig.MaxAge
			}
			if corsConfig.AllowCredentials != nil {
				output["allow_credentials"] = *corsConfig.AllowCredentials
			}
			corsOutput = append(corsOutput, output)
		}
		/*		if err := d.Set("cors_configuration", corsOutput); err != nil {
				return fmt.Errorf("Error setting `cors_configuration`: %+v", corsOutput)
			}*/
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmHealthcareServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).healthcare.HealthcareServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error Parsing Azure Resource ID: %+v", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["services"]
	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deleting Healthcare Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
