// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	service "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultSuppress "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/suppress"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
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
			_, err := service.ParseServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(service.KindFhir),
				ValidateFunc: validation.StringInSlice([]string{
					string(service.KindFhir),
					string(service.KindFhirNegativeRFour),
					string(service.KindFhirNegativeStuThree),
				}, false),
			},

			"cosmosdb_throughput": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1000,
				ValidateFunc: validation.IntBetween(1, 100000),
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
									"PATCH",
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := service.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.ServicesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_healthcare_service", id.ID())
		}
	}

	cosmosDbConfiguration, err := expandsCosmosDBConfiguration(d)
	if err != nil {
		return fmt.Errorf("expanding cosmosdb_configuration: %+v", err)
	}

	healthcareServiceDescription := service.ServicesDescription{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Kind:     service.Kind(d.Get("kind").(string)),
		Properties: &service.ServicesProperties{
			AccessPolicies:              expandAccessPolicyEntries(d),
			CosmosDbConfiguration:       cosmosDbConfiguration,
			CorsConfiguration:           expandCorsConfiguration(d),
			AuthenticationConfiguration: expandAuthentication(d),
		},
	}

	publicNetworkAccess := d.Get("public_network_access_enabled").(bool)
	if !publicNetworkAccess {
		healthcareServiceDescription.Properties.PublicNetworkAccess = pointer.To(service.PublicNetworkAccessDisabled)
	} else {
		healthcareServiceDescription.Properties.PublicNetworkAccess = pointer.To(service.PublicNetworkAccessEnabled)
	}

	err = client.ServicesCreateOrUpdateThenPoll(ctx, id, healthcareServiceDescription)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHealthcareServiceRead(d, meta)
}

func resourceHealthcareServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := service.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ServicesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if m := resp.Model; m != nil {
		d.Set("location", location.Normalize(m.Location))

		if kind := m.Kind; string(kind) != "" {
			d.Set("kind", kind)
		}
		if props := m.Properties; props != nil {
			if err := d.Set("access_policy_object_ids", flattenAccessPolicies(props.AccessPolicies)); err != nil {
				return fmt.Errorf("setting `access_policy_object_ids`: %+v", err)
			}

			cosmodDbKeyVaultKeyVersionlessId := ""
			cosmosDbThroughput := 0
			if cosmos := props.CosmosDbConfiguration; cosmos != nil {
				if v := cosmos.OfferThroughput; v != nil {
					cosmosDbThroughput = int(*v)
				}
				if v := cosmos.KeyVaultKeyUri; v != nil {
					cosmodDbKeyVaultKeyVersionlessId = *v
				}
			}
			d.Set("cosmosdb_key_vault_key_versionless_id", cosmodDbKeyVaultKeyVersionlessId)
			d.Set("cosmosdb_throughput", cosmosDbThroughput)
			if pointer.From(props.PublicNetworkAccess) == service.PublicNetworkAccessEnabled {
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

		return tags.FlattenAndSet(d, m.Tags)
	}

	return nil
}

func resourceHealthcareServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := service.ParseServiceID(d.Id())
	if err != nil {
		return fmt.Errorf("Parsing Azure Resource ID: %+v", err)
	}

	err = client.ServicesDeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Healthcare Service %q (Resource Group %q): %+v", id.ServiceName, id.ResourceGroupName, err)
	}
	return nil
}

func expandAccessPolicyEntries(d *pluginsdk.ResourceData) *[]service.ServiceAccessPolicyEntry {
	accessPolicyObjectIds := d.Get("access_policy_object_ids").(*pluginsdk.Set).List()
	svcAccessPolicyArray := make([]service.ServiceAccessPolicyEntry, 0)

	for _, objectId := range accessPolicyObjectIds {
		svcAccessPolicyObjectId := service.ServiceAccessPolicyEntry{ObjectId: objectId.(string)}
		svcAccessPolicyArray = append(svcAccessPolicyArray, svcAccessPolicyObjectId)
	}

	return &svcAccessPolicyArray
}

func expandCorsConfiguration(d *pluginsdk.ResourceData) *service.ServiceCorsConfigurationInfo {
	corsConfigRaw := d.Get("cors_configuration").([]interface{})

	if len(corsConfigRaw) == 0 {
		return &service.ServiceCorsConfigurationInfo{}
	}

	corsConfigAttr := corsConfigRaw[0].(map[string]interface{})

	allowedOrigins := *utils.ExpandStringSlice(corsConfigAttr["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *utils.ExpandStringSlice(corsConfigAttr["allowed_headers"].(*pluginsdk.Set).List())
	allowedMethods := *utils.ExpandStringSlice(corsConfigAttr["allowed_methods"].([]interface{}))
	maxAgeInSeconds := int64(corsConfigAttr["max_age_in_seconds"].(int))
	allowCredentials := corsConfigAttr["allow_credentials"].(bool)

	cors := &service.ServiceCorsConfigurationInfo{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		MaxAge:           &maxAgeInSeconds,
		AllowCredentials: &allowCredentials,
	}
	return cors
}

func expandAuthentication(d *pluginsdk.ResourceData) *service.ServiceAuthenticationConfigurationInfo {
	authConfigRaw := d.Get("authentication_configuration").([]interface{})

	if len(authConfigRaw) == 0 {
		return &service.ServiceAuthenticationConfigurationInfo{}
	}

	authConfigAttr := authConfigRaw[0].(map[string]interface{})
	authority := authConfigAttr["authority"].(string)
	audience := authConfigAttr["audience"].(string)
	smartProxyEnabled := authConfigAttr["smart_proxy_enabled"].(bool)

	auth := &service.ServiceAuthenticationConfigurationInfo{
		Authority:         &authority,
		Audience:          &audience,
		SmartProxyEnabled: &smartProxyEnabled,
	}
	return auth
}

func expandsCosmosDBConfiguration(d *pluginsdk.ResourceData) (*service.ServiceCosmosDbConfigurationInfo, error) {
	cosmosdb := &service.ServiceCosmosDbConfigurationInfo{
		OfferThroughput: pointer.To(int64(d.Get("cosmosdb_throughput").(int))),
	}

	if keyVaultKeyIDRaw, ok := d.GetOk("cosmosdb_key_vault_key_versionless_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyIDRaw.(string))
		if err != nil {
			return nil, fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}
		cosmosdb.KeyVaultKeyUri = pointer.To(keyVaultKey.ID())
	}

	return cosmosdb, nil
}

func flattenAccessPolicies(policies *[]service.ServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		result = append(result, policy.ObjectId)

	}

	return result
}

func flattenAuthentication(input *service.ServiceAuthenticationConfigurationInfo) []interface{} {
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

func flattenCorsConfig(input *service.ServiceCorsConfigurationInfo) []interface{} {
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
