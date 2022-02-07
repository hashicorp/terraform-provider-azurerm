package healthcare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2020-03-30/healthcareapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultSuppress "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/suppress"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareServiceCreateUpdate,
		Read:   resourceHealthcareServiceRead,
		Update: resourceHealthcareServiceCreateUpdate,
		Delete: resourceHealthcareServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(healthcareapis.Fhir),
				ValidateFunc: validation.StringInSlice([]string{
					string(healthcareapis.Fhir),
					string(healthcareapis.FhirR4),
					string(healthcareapis.FhirStu3),
				}, false),
			},

			"cosmosdb_throughput": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1000,
				ValidateFunc: validation.IntBetween(1, 10000),
			},

			"cosmosdb_key_vault_key_versionless_id": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: keyVaultSuppress.DiffSuppressIgnoreKeyVaultKeyVersion,
				ValidateFunc:     keyVaultValidate.VersionlessNestedItemId,
			},

			"access_policy_object_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"authentication_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"authority": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
						},
						"audience": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
						},
						"smart_proxy_enabled": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
						},
					},
				},
			},

			"cors_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{
								"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
								"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
								"cors_configuration.0.allow_credentials",
							},
						},
						"allowed_headers": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{
								"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
								"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
								"cors_configuration.0.allow_credentials",
							},
						},
						"allowed_methods": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 64,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
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
							AtLeastOneOf: []string{
								"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
								"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
								"cors_configuration.0.allow_credentials",
							},
						},
						"max_age_in_seconds": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 2000000000),
							AtLeastOneOf: []string{
								"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
								"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
								"cors_configuration.0.allow_credentials",
							},
						},
						"allow_credentials": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
								"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
								"cors_configuration.0.allow_credentials",
							},
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceHealthcareServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_healthcare_service", id.ID())
		}
	}

	cosmosDbConfiguration, err := expandsCosmosDBConfiguration(d)
	if err != nil {
		return fmt.Errorf("expanding cosmosdb_configuration: %+v", err)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	kind := d.Get("kind").(string)

	healthcareServiceDescription := healthcareapis.ServicesDescription{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		Kind:     healthcareapis.Kind(kind),
		Properties: &healthcareapis.ServicesProperties{
			AccessPolicies:              expandAccessPolicyEntries(d),
			CosmosDbConfiguration:       cosmosDbConfiguration,
			CorsConfiguration:           expandCorsConfiguration(d),
			AuthenticationConfiguration: expandAuthentication(d),
		},
	}

	publicNetworkAccess := d.Get("public_network_access_enabled").(bool)
	if !publicNetworkAccess {
		healthcareServiceDescription.Properties.PublicNetworkAccess = healthcareapis.Disabled
	} else {
		healthcareServiceDescription.Properties.PublicNetworkAccess = healthcareapis.Enabled
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, healthcareServiceDescription)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareServiceRead(d, meta)
}

func resourceHealthcareServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[WARN] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
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
		if err := d.Set("access_policy_object_ids", flattenAccessPolicies(props.AccessPolicies)); err != nil {
			return fmt.Errorf("setting `access_policy_object_ids`: %+v", err)
		}

		cosmodDbKeyVaultKeyVersionlessId := ""
		cosmosDbThroughput := 0
		if cosmos := props.CosmosDbConfiguration; cosmos != nil {
			if v := cosmos.OfferThroughput; v != nil {
				cosmosDbThroughput = int(*v)
			}
			if v := cosmos.KeyVaultKeyURI; v != nil {
				cosmodDbKeyVaultKeyVersionlessId = *v
			}
		}
		d.Set("cosmosdb_key_vault_key_versionless_id", cosmodDbKeyVaultKeyVersionlessId)
		d.Set("cosmosdb_throughput", cosmosDbThroughput)
		if props.PublicNetworkAccess == healthcareapis.Enabled {
			d.Set("public_network_access_enabled", true)
		} else {
			d.Set("public_network_access_enabled", false)
		}

		if err := d.Set("authentication_configuration", flattenAuthentication(props.AuthenticationConfiguration)); err != nil {
			return fmt.Errorf("setting `authentication_configuration`: %+v", err)
		}

		if err := d.Set("cors_configuration", flattenCorsConfig(props.CorsConfiguration)); err != nil {
			return fmt.Errorf("setting `cors_configuration`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceHealthcareServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("Parsing Azure Resource ID: %+v", err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Healthcare Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deleting Healthcare Service %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandAccessPolicyEntries(d *pluginsdk.ResourceData) *[]healthcareapis.ServiceAccessPolicyEntry {
	accessPolicyObjectIds := d.Get("access_policy_object_ids").(*pluginsdk.Set).List()
	svcAccessPolicyArray := make([]healthcareapis.ServiceAccessPolicyEntry, 0)

	for _, objectId := range accessPolicyObjectIds {
		svcAccessPolicyObjectId := healthcareapis.ServiceAccessPolicyEntry{ObjectID: utils.String(objectId.(string))}
		svcAccessPolicyArray = append(svcAccessPolicyArray, svcAccessPolicyObjectId)
	}

	return &svcAccessPolicyArray
}

func expandCorsConfiguration(d *pluginsdk.ResourceData) *healthcareapis.ServiceCorsConfigurationInfo {
	corsConfigRaw := d.Get("cors_configuration").([]interface{})

	if len(corsConfigRaw) == 0 {
		return &healthcareapis.ServiceCorsConfigurationInfo{}
	}

	corsConfigAttr := corsConfigRaw[0].(map[string]interface{})

	allowedOrigins := *utils.ExpandStringSlice(corsConfigAttr["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *utils.ExpandStringSlice(corsConfigAttr["allowed_headers"].(*pluginsdk.Set).List())
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

func expandAuthentication(d *pluginsdk.ResourceData) *healthcareapis.ServiceAuthenticationConfigurationInfo {
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

func expandsCosmosDBConfiguration(d *pluginsdk.ResourceData) (*healthcareapis.ServiceCosmosDbConfigurationInfo, error) {
	throughput := int32(d.Get("cosmosdb_throughput").(int))

	cosmosdb := &healthcareapis.ServiceCosmosDbConfigurationInfo{
		OfferThroughput: utils.Int32(throughput),
	}

	if keyVaultKeyIDRaw, ok := d.GetOk("cosmosdb_key_vault_key_versionless_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyIDRaw.(string))
		if err != nil {
			return nil, fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}
		cosmosdb.KeyVaultKeyURI = utils.String(keyVaultKey.ID())
	}

	return cosmosdb, nil
}

func flattenAccessPolicies(policies *[]healthcareapis.ServiceAccessPolicyEntry) []string {
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

func flattenAuthentication(input *healthcareapis.ServiceAuthenticationConfigurationInfo) []interface{} {
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

func flattenCorsConfig(input *healthcareapis.ServiceCorsConfigurationInfo) []interface{} {
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
