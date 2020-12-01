package healthcare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2019-09-16/healthcareapis"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/healthcare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHealthcareService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHealthcareServiceCreateUpdate,
		Read:   resourceArmHealthcareServiceRead,
		Update: resourceArmHealthcareServiceCreateUpdate,
		Delete: resourceArmHealthcareServiceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServiceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"authentication_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"allowed_headers": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 64,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"DELETE",
									"GET",
									"HEAD",
									"MERGE",
									"POST",
									"OPTIONS",
									"PUT",
								}, false),
							},
						},
						"max_age_in_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 2000000000),
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
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Healthcare Service creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	kind := d.Get("kind").(string)
	cdba := int32(d.Get("cosmosdb_throughput").(int))

	if d.IsNewResource() {
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
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Healthcare Service %q was not found (Resource Group %q)", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Healthcare Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if kind := resp.Kind; string(kind) != "" {
		d.Set("kind", kind)
	}
	if props := resp.Properties; props != nil {
		if err := d.Set("access_policy_object_ids", flattenHealthcareAccessPolicies(props.AccessPolicies)); err != nil {
			return fmt.Errorf("Error setting `access_policy_object_ids`: %+v", err)
		}

		cosmosThroughput := 0
		if props.CosmosDbConfiguration != nil && props.CosmosDbConfiguration.OfferThroughput != nil {
			cosmosThroughput = int(*props.CosmosDbConfiguration.OfferThroughput)
		}
		d.Set("cosmosdb_throughput", cosmosThroughput)

		if err := d.Set("authentication_configuration", flattenHealthcareAuthConfig(props.AuthenticationConfiguration)); err != nil {
			return fmt.Errorf("Error setting `authentication_configuration`: %+v", err)
		}

		if err := d.Set("cors_configuration", flattenHealthcareCorsConfig(props.CorsConfiguration)); err != nil {
			return fmt.Errorf("Error setting `cors_configuration`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmHealthcareServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error Parsing Azure Resource ID: %+v", err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Healthcare Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deleting Healthcare Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
	smartProxyEnabled := authConfigAttr["smart_proxy_enabled"].(bool)

	auth := &healthcareapis.ServiceAuthenticationConfigurationInfo{
		Authority:         &authority,
		Audience:          &audience,
		SmartProxyEnabled: &smartProxyEnabled,
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

func flattenHealthcareAuthConfig(input *healthcareapis.ServiceAuthenticationConfigurationInfo) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	authority := ""
	if input.Authority != nil {
		authority = *input.Authority
	}
	audience := ""
	if input.Audience != nil {
		audience = *input.Audience
	}
	smartProxyEnabled := false
	if input.SmartProxyEnabled != nil {
		smartProxyEnabled = *input.SmartProxyEnabled
	}
	return []interface{}{
		map[string]interface{}{
			"audience":            audience,
			"authority":           authority,
			"smart_proxy_enabled": smartProxyEnabled,
		},
	}
}

func flattenHealthcareCorsConfig(input *healthcareapis.ServiceCorsConfigurationInfo) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	maxAge := 0
	if input.MaxAge != nil {
		maxAge = int(*input.MaxAge)
	}
	allowCredentials := false
	if input.AllowCredentials != nil {
		allowCredentials = *input.AllowCredentials
	}

	return []interface{}{
		map[string]interface{}{
			"allow_credentials":  allowCredentials,
			"allowed_headers":    utils.FlattenStringSlice(input.Headers),
			"allowed_methods":    utils.FlattenStringSlice(input.Methods),
			"allowed_origins":    utils.FlattenStringSlice(input.Origins),
			"max_age_in_seconds": maxAge,
		},
	}
}
