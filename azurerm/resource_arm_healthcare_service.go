package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2019-09-16/healthcareapis"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(healthcareapis.Fhir),
				ValidateFunc: validation.StringInSlice([]string{
					string(healthcareapis.Fhir),
					string(healthcareapis.FhirR4),
					string(healthcareapis.FhirStu3),
				}, false),
			},

			"cosmosdb_throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1000,
				ValidateFunc: validation.IntBetween(1, 10000),
			},

			"access_policy_object_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.UUID,
				},
			},

			"authentication_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authority": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"audience": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smart_proxy_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"cors_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"allowed_headers": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								Elem: &schema.Schema{
									Type: schema.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										"DELETE",
										"GET",
										"HEAD",
										"MERGE",
										"POST",
										"OPTIONS",
										"PUT"}, false),
								},
							},
						},
						"max_age_in_seconds": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 2000000000),
						},
						"allow_credentials": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmHealthcareServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Healthcare Service creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	kind := d.Get("kind").(string)
	cdba := int32(d.Get("cosmosdb_throughput").(int))

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	healthcareServiceDescription := healthcareapis.ServicesDescription{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		Kind:     healthcareapis.Kind(kind),
		Properties: &healthcareapis.ServicesProperties{
			AccessPolicies: expandAzureRMhealthcareapisAccessPolicyEntries(d),
			CosmosDbConfiguration: &healthcareapis.ServiceCosmosDbConfigurationInfo{
				OfferThroughput: &cdba,
			},
			CorsConfiguration:           expandAzureRMhealthcareapisCorsConfiguration(d),
			AuthenticationConfiguration: expandAzureRMhealthcareapisAuthentication(d),
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
	client := meta.(*ArmClient).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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

	if kind := resp.Kind; string(kind) != "" {
		d.Set("kind", kind)
	}
	if properties := resp.Properties; properties != nil {
		if config := properties.AccessPolicies; config != nil {
			d.Set("access_policy_object_ids", flattenHealthcareAccessPolicies(config))
		}
		if config := properties.CosmosDbConfiguration; config != nil {
			d.Set("cosmosdb_throughput", config.OfferThroughput)
		}

		if authConfig := properties.AuthenticationConfiguration; authConfig != nil {
			if err := d.Set("authentication_configuration", flattenHealthcareAuthConfig(authConfig)); err != nil {
				return fmt.Errorf("Error setting `authentication_configuration`: %+v", flattenHealthcareAuthConfig(authConfig))
			}
		}

		if corsConfig := properties.CorsConfiguration; corsConfig != nil {
			if err := d.Set("cors_configuration", flattenHealthcareCorsConfig(corsConfig)); err != nil {
				return fmt.Errorf("Error setting `cors_configuration`: %+v", flattenHealthcareCorsConfig(corsConfig))
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmHealthcareServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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

func expandAzureRMhealthcareapisAccessPolicyEntries(d *schema.ResourceData) *[]healthcareapis.ServiceAccessPolicyEntry {
	accessPolicyObjectIds := d.Get("access_policy_object_ids").(*schema.Set).List()
	svcAccessPolicyArray := make([]healthcareapis.ServiceAccessPolicyEntry, 0)

	for _, objectId := range accessPolicyObjectIds {
		svcAccessPolicyObjectId := healthcareapis.ServiceAccessPolicyEntry{ObjectID: utils.String(objectId.(string))}
		svcAccessPolicyArray = append(svcAccessPolicyArray, svcAccessPolicyObjectId)
	}

	return &svcAccessPolicyArray
}

func expandAzureRMhealthcareapisCorsConfiguration(d *schema.ResourceData) *healthcareapis.ServiceCorsConfigurationInfo {
	corsConfigRaw := d.Get("cors_configuration").([]interface{})

	if len(corsConfigRaw) == 0 {
		return &healthcareapis.ServiceCorsConfigurationInfo{}
	}

	corsConfigAttr := corsConfigRaw[0].(map[string]interface{})

	allowedOrigins := *utils.ExpandStringSlice(corsConfigAttr["allowed_origins"].(*schema.Set).List())
	allowedHeaders := *utils.ExpandStringSlice(corsConfigAttr["allowed_headers"].(*schema.Set).List())
	allowedMethods := *utils.ExpandStringSlice(corsConfigAttr["allowed_methods"].([]interface{}))
	maxAgeInSeconds := int32(corsConfigAttr["max_age_in_seconds"].(int))
	allowCredentials := corsConfigAttr["allow_credentials"].(bool)

	cors := &healthcareapis.ServiceCorsConfigurationInfo{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		MaxAge:           &maxAgeInSeconds,
		AllowCredentials: &allowCredentials,
	}
	return cors
}

func expandAzureRMhealthcareapisAuthentication(d *schema.ResourceData) *healthcareapis.ServiceAuthenticationConfigurationInfo {
	authConfigRaw := d.Get("authentication_configuration").([]interface{})

	if len(authConfigRaw) == 0 {
		return &healthcareapis.ServiceAuthenticationConfigurationInfo{}
	}

	authConfigAttr := authConfigRaw[0].(map[string]interface{})
	authority := authConfigAttr["authority"].(string)
	audience := authConfigAttr["audience"].(string)
	smart_proxy_enabled := authConfigAttr["smart_proxy_enabled"].(bool)

	auth := &healthcareapis.ServiceAuthenticationConfigurationInfo{
		Authority:         &authority,
		Audience:          &audience,
		SmartProxyEnabled: &smart_proxy_enabled,
	}
	return auth
}

func flattenHealthcareAccessPolicies(policies *[]healthcareapis.ServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		if objectId := policy.ObjectID; objectId != nil {
			result = append(result, *objectId)
		}
	}

	return result
}

func flattenHealthcareAuthConfig(authConfig *healthcareapis.ServiceAuthenticationConfigurationInfo) []interface{} {
	authOutput := make([]interface{}, 0)

	output := make(map[string]interface{})
	if authConfig.Authority != nil {
		output["authority"] = *authConfig.Authority
	}
	if authConfig.Audience != nil {
		output["audience"] = *authConfig.Audience
	}
	if authConfig.SmartProxyEnabled != nil {
		output["smart_proxy_enabled"] = *authConfig.SmartProxyEnabled
	}
	authOutput = append(authOutput, output)

	return authOutput
}

func flattenHealthcareCorsConfig(corsConfig *healthcareapis.ServiceCorsConfigurationInfo) []interface{} {
	corsOutput := make([]interface{}, 0)

	output := make(map[string]interface{})
	if corsConfig.Origins != nil {
		output["allowed_origins"] = *corsConfig.Origins
	}
	if corsConfig.Headers != nil {
		output["allowed_headers"] = *corsConfig.Headers
	}
	if corsConfig.Methods != nil {
		output["allowed_methods"] = *corsConfig.Methods
	}
	if corsConfig.MaxAge != nil {
		output["max_age_in_seconds"] = *corsConfig.MaxAge
	}
	if corsConfig.AllowCredentials != nil {
		output["allow_credentials"] = *corsConfig.AllowCredentials
	}
	corsOutput = append(corsOutput, output)

	return corsOutput
}
