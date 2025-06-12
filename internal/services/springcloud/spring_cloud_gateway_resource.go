// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudGatewayModel struct {
	Name                                  string                     `tfschema:"name"`
	SpringCloudServiceId                  string                     `tfschema:"spring_cloud_service_id"`
	ApiMetadata                           []ApiMetadataModel         `tfschema:"api_metadata"`
	ApplicationPerformanceMonitoringIds   []string                   `tfschema:"application_performance_monitoring_ids"`
	ApplicationPerformanceMonitoringTypes []string                   `tfschema:"application_performance_monitoring_types"`
	ClientAuthorization                   []ClientAuthorizationModel `tfschema:"client_authorization"`
	Cors                                  []CorsModel                `tfschema:"cors"`
	EnvironmentVariables                  map[string]string          `tfschema:"environment_variables"`
	LocalResponseCachePerRoute            []ResponseCacheModel       `tfschema:"local_response_cache_per_route"`
	LocalResponseCachePerInstance         []ResponseCacheModel       `tfschema:"local_response_cache_per_instance"`
	SensitiveEnvironmentVariables         map[string]string          `tfschema:"sensitive_environment_variables"`
	HttpsOnly                             bool                       `tfschema:"https_only"`
	InstanceCount                         int64                      `tfschema:"instance_count"`
	PublicNetworkAccessEnabled            bool                       `tfschema:"public_network_access_enabled"`
	Quota                                 []QuotaModel               `tfschema:"quota"`
	Sso                                   []GatewaySsoModel          `tfschema:"sso"`
	Url                                   string                     `tfschema:"url"`
}

type ApiMetadataModel struct {
	Description      string `tfschema:"description"`
	DocumentationUrl string `tfschema:"documentation_url"`
	ServerURL        string `tfschema:"server_url"`
	Title            string `tfschema:"title"`
	Version          string `tfschema:"version"`
}

type ClientAuthorizationModel struct {
	CertificateIds      []string `tfschema:"certificate_ids"`
	VerificationEnabled bool     `tfschema:"verification_enabled"`
}

type CorsModel struct {
	CredentialsAllowed    bool     `tfschema:"credentials_allowed"`
	AllowedHeaders        []string `tfschema:"allowed_headers"`
	AllowedMethods        []string `tfschema:"allowed_methods"`
	AllowedOrigins        []string `tfschema:"allowed_origins"`
	AllowedOriginPatterns []string `tfschema:"allowed_origin_patterns"`
	ExposedHeaders        []string `tfschema:"exposed_headers"`
	MaxAgeSeconds         int64    `tfschema:"max_age_seconds"`
}

type GatewaySsoModel struct {
	ClientId     string   `tfschema:"client_id"`
	ClientSecret string   `tfschema:"client_secret"`
	IssuerUri    string   `tfschema:"issuer_uri"`
	Scope        []string `tfschema:"scope"`
}

type QuotaModel struct {
	Cpu    string `tfschema:"cpu"`
	Memory string `tfschema:"memory"`
}

type ResponseCacheModel struct {
	Size       string `tfschema:"size"`
	TimeToLive string `tfschema:"time_to_live"`
}

type SpringCloudGatewayResource struct{}

func (s SpringCloudGatewayResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_gateway` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudGatewayResource{}
	_ sdk.ResourceWithStateMigration              = SpringCloudGatewayResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudGatewayResource{}
)

func (s SpringCloudGatewayResource) ResourceType() string {
	return "azurerm_spring_cloud_gateway"
}

func (s SpringCloudGatewayResource) ModelObject() interface{} {
	return &SpringCloudGatewayModel{}
}

func (s SpringCloudGatewayResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateGatewayID
}

func (s SpringCloudGatewayResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"default",
			}, false),
		},

		"spring_cloud_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSpringCloudServiceID,
		},

		"api_metadata": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"documentation_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"server_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"title": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"application_performance_monitoring_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: appplatform.ValidateApmID,
			},
		},

		"application_performance_monitoring_types": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(appplatform.ApmTypeAppDynamics),
					string(appplatform.ApmTypeApplicationInsights),
					string(appplatform.ApmTypeDynatrace),
					string(appplatform.ApmTypeElasticAPM),
					string(appplatform.ApmTypeNewRelic),
				}, false),
			},
		},

		"client_authorization": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: appplatform.ValidateCertificateID,
						},
					},

					"verification_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"credentials_allowed": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"allowed_headers": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_methods": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
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
					},

					"allowed_origins": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"allowed_origin_patterns": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"exposed_headers": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"max_age_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"environment_variables": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"local_response_cache_per_route": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"local_response_cache_per_instance"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"size": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"time_to_live": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"local_response_cache_per_instance": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"local_response_cache_per_route"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"size": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"time_to_live": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"sensitive_environment_variables": {
			Type:      pluginsdk.TypeMap,
			Optional:  true,
			Sensitive: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 500),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"quota": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "1",
						// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"memory": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "2Gi",
						// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"sso": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"client_secret": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"issuer_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"scope": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func (s SpringCloudGatewayResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (s SpringCloudGatewayResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudGatewayV0ToV1{},
		},
	}
}

func (s SpringCloudGatewayResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := appplatform.NewGatewayID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			existing, err := client.GatewaysGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			service, err := client.ServicesGet(ctx, *springId)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", springId, err)
			}
			if service.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", springId)
			}
			if service.Model.Sku == nil || service.Model.Sku.Name == nil || service.Model.Sku.Tier == nil {
				return fmt.Errorf("invalid `sku` for %s", springId)
			}

			gatewayResource := appplatform.GatewayResource{
				Properties: &appplatform.GatewayProperties{
					ClientAuth:              expandGatewayClientAuth(model.ClientAuthorization),
					ApiMetadataProperties:   expandGatewayGatewayAPIMetadataProperties(model.ApiMetadata),
					Apms:                    expandGatewayApms(model.ApplicationPerformanceMonitoringIds),
					ApmTypes:                expandGatewayGatewayApmTypes(model.ApplicationPerformanceMonitoringTypes),
					CorsProperties:          expandGatewayGatewayCorsProperties(model.Cors),
					EnvironmentVariables:    expandGatewayGatewayEnvironmentVariables(model.EnvironmentVariables, model.SensitiveEnvironmentVariables),
					HTTPSOnly:               pointer.To(model.HttpsOnly),
					Public:                  pointer.To(model.PublicNetworkAccessEnabled),
					ResponseCacheProperties: expandGatewayResponseCacheProperties(model),
					ResourceRequests:        expandGatewayGatewayResourceRequests(model.Quota),
					SsoProperties:           expandGatewaySsoProperties(model.Sso),
				},
				Sku: &appplatform.Sku{
					Name:     service.Model.Sku.Name,
					Tier:     service.Model.Sku.Tier,
					Capacity: pointer.To(model.InstanceCount),
				},
			}

			err = client.GatewaysCreateOrUpdateThenPoll(ctx, id, gatewayResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudGatewayResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := appplatform.ParseGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			resp, err := client.GatewaysGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			properties := resp.Model.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			sku := resp.Model.Sku
			if sku == nil {
				return fmt.Errorf("retrieving %s: sku was nil", id)
			}

			if metadata.ResourceData.HasChange("client_authorization") {
				properties.ClientAuth = expandGatewayClientAuth(model.ClientAuthorization)
			}

			if metadata.ResourceData.HasChange("api_metadata") {
				properties.ApiMetadataProperties = expandGatewayGatewayAPIMetadataProperties(model.ApiMetadata)
			}

			if metadata.ResourceData.HasChange("application_performance_monitoring_ids") {
				properties.Apms = expandGatewayApms(model.ApplicationPerformanceMonitoringIds)
			}

			if metadata.ResourceData.HasChange("application_performance_monitoring_types") {
				properties.ApmTypes = expandGatewayGatewayApmTypes(model.ApplicationPerformanceMonitoringTypes)
			}

			if metadata.ResourceData.HasChange("cors") {
				properties.CorsProperties = expandGatewayGatewayCorsProperties(model.Cors)
			}

			if metadata.ResourceData.HasChange("environment_variables") || metadata.ResourceData.HasChange("sensitive_environment_variables") {
				properties.EnvironmentVariables = expandGatewayGatewayEnvironmentVariables(model.EnvironmentVariables, model.SensitiveEnvironmentVariables)
			}

			if metadata.ResourceData.HasChange("https_only") {
				properties.HTTPSOnly = pointer.To(model.HttpsOnly)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				properties.Public = pointer.To(model.PublicNetworkAccessEnabled)
			}

			if metadata.ResourceData.HasChange("quota") {
				properties.ResourceRequests = expandGatewayGatewayResourceRequests(model.Quota)
			}

			if metadata.ResourceData.HasChange("sso") {
				properties.SsoProperties = expandGatewaySsoProperties(model.Sso)
			}

			if metadata.ResourceData.HasChange("local_response_cache_per_instance") || metadata.ResourceData.HasChange("local_response_cache_per_route") {
				properties.ResponseCacheProperties = expandGatewayResponseCacheProperties(model)
			}

			if metadata.ResourceData.HasChange("instance_count") {
				sku.Capacity = pointer.To(model.InstanceCount)
			}
			resource := appplatform.GatewayResource{
				Properties: properties,
				Sku:        sku,
			}

			err = client.GatewaysCreateOrUpdateThenPoll(ctx, *id, resource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s SpringCloudGatewayResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GatewaysGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)

			var model SpringCloudGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudGatewayModel{
				Name:                          id.GatewayName,
				SpringCloudServiceId:          springId.ID(),
				SensitiveEnvironmentVariables: model.SensitiveEnvironmentVariables,
			}

			if resp.Model != nil {
				if props := resp.Model.Properties; props != nil {
					state.ApiMetadata = flattenGatewayGatewayAPIMetadataProperties(props.ApiMetadataProperties)
					apms, err := flattenGatewayApms(props.Apms)
					if err != nil {
						return err
					}
					state.ApplicationPerformanceMonitoringIds = apms
					state.ApplicationPerformanceMonitoringTypes = flattenGatewayGatewayApmTypes(props.ApmTypes)
					state.ClientAuthorization = flattenGatewayClientAuth(props.ClientAuth)
					state.Cors = flattenGatewayGatewayCorsProperties(props.CorsProperties)
					if props.EnvironmentVariables != nil {
						state.EnvironmentVariables = pointer.From(props.EnvironmentVariables.Properties)
					}
					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.PublicNetworkAccessEnabled = pointer.From(props.Public)
					state.Quota = flattenGatewayGatewayResourceRequests(props.ResourceRequests)
					state.Sso = flattenGatewaySsoProperties(props.SsoProperties, model.Sso)
					state.LocalResponseCachePerRoute = flattenGatewayLocalResponseCachePerRouteProperties(props.ResponseCacheProperties)
					state.LocalResponseCachePerInstance = flattenGatewayLocalResponseCachePerInstanceProperties(props.ResponseCacheProperties)
				}

				if sku := resp.Model.Sku; sku != nil {
					state.InstanceCount = pointer.From(sku.Capacity)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudGatewayResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.GatewaysDeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandGatewayGatewayAPIMetadataProperties(input []ApiMetadataModel) *appplatform.GatewayApiMetadataProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.GatewayApiMetadataProperties{
		Title:         pointer.To(v.Title),
		Description:   pointer.To(v.Description),
		Documentation: pointer.To(v.DocumentationUrl),
		Version:       pointer.To(v.Version),
		ServerURL:     pointer.To(v.ServerURL),
	}
}

func expandGatewayGatewayCorsProperties(input []CorsModel) *appplatform.GatewayCorsProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.GatewayCorsProperties{
		AllowedOrigins:        pointer.To(v.AllowedOrigins),
		AllowedOriginPatterns: pointer.To(v.AllowedOriginPatterns),
		AllowedMethods:        pointer.To(v.AllowedMethods),
		AllowedHeaders:        pointer.To(v.AllowedHeaders),
		MaxAge:                pointer.To(v.MaxAgeSeconds),
		AllowCredentials:      pointer.To(v.CredentialsAllowed),
		ExposedHeaders:        pointer.To(v.ExposedHeaders),
	}
}

func expandGatewayGatewayResourceRequests(input []QuotaModel) *appplatform.GatewayResourceRequests {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.GatewayResourceRequests{
		Cpu:    pointer.To(v.Cpu),
		Memory: pointer.To(v.Memory),
	}
}

func expandGatewaySsoProperties(input []GatewaySsoModel) *appplatform.SsoProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.SsoProperties{
		Scope:        pointer.To(v.Scope),
		ClientId:     pointer.To(v.ClientId),
		ClientSecret: pointer.To(v.ClientSecret),
		IssuerUri:    pointer.To(v.IssuerUri),
	}
}

func expandGatewayGatewayApmTypes(input []string) *[]appplatform.ApmType {
	if len(input) == 0 {
		return nil
	}
	out := make([]appplatform.ApmType, 0)
	for _, v := range input {
		out = append(out, appplatform.ApmType(v))
	}
	return &out
}

func expandGatewayGatewayEnvironmentVariables(env map[string]string, secrets map[string]string) *appplatform.GatewayPropertiesEnvironmentVariables {
	if len(env) == 0 && len(secrets) == 0 {
		return nil
	}

	return &appplatform.GatewayPropertiesEnvironmentVariables{
		Properties: pointer.To(env),
		Secrets:    pointer.To(secrets),
	}
}

func expandGatewayClientAuth(input []ClientAuthorizationModel) *appplatform.GatewayPropertiesClientAuth {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	verificationEnabled := appplatform.GatewayCertificateVerificationDisabled
	if v.VerificationEnabled {
		verificationEnabled = appplatform.GatewayCertificateVerificationEnabled
	}
	return &appplatform.GatewayPropertiesClientAuth{
		Certificates:            pointer.To(v.CertificateIds),
		CertificateVerification: pointer.To(verificationEnabled),
	}
}

func expandGatewayResponseCacheProperties(input SpringCloudGatewayModel) appplatform.GatewayResponseCacheProperties {
	if len(input.LocalResponseCachePerRoute) != 0 {
		return appplatform.GatewayLocalResponseCachePerRouteProperties{
			Size:       pointer.To(input.LocalResponseCachePerRoute[0].Size),
			TimeToLive: pointer.To(input.LocalResponseCachePerRoute[0].TimeToLive),
		}
	}
	if len(input.LocalResponseCachePerInstance) != 0 {
		return appplatform.GatewayLocalResponseCachePerInstanceProperties{
			Size:       pointer.To(input.LocalResponseCachePerInstance[0].Size),
			TimeToLive: pointer.To(input.LocalResponseCachePerInstance[0].TimeToLive),
		}
	}
	return appplatform.RawGatewayResponseCachePropertiesImpl{}
}

func expandGatewayApms(input []string) *[]appplatform.ApmReference {
	if len(input) == 0 {
		return nil
	}
	out := make([]appplatform.ApmReference, 0)
	for _, v := range input {
		out = append(out, appplatform.ApmReference{
			ResourceId: v,
		})
	}
	return pointer.To(out)
}

func flattenGatewayGatewayAPIMetadataProperties(input *appplatform.GatewayApiMetadataProperties) []ApiMetadataModel {
	if input == nil {
		return make([]ApiMetadataModel, 0)
	}

	return []ApiMetadataModel{
		{
			Description:      pointer.From(input.Description),
			DocumentationUrl: pointer.From(input.Documentation),
			ServerURL:        pointer.From(input.ServerURL),
			Title:            pointer.From(input.Title),
			Version:          pointer.From(input.Version),
		},
	}
}

func flattenGatewayGatewayCorsProperties(input *appplatform.GatewayCorsProperties) []CorsModel {
	if input == nil {
		return make([]CorsModel, 0)
	}

	return []CorsModel{
		{
			CredentialsAllowed:    pointer.From(input.AllowCredentials),
			AllowedHeaders:        pointer.From(input.AllowedHeaders),
			AllowedMethods:        pointer.From(input.AllowedMethods),
			AllowedOrigins:        pointer.From(input.AllowedOrigins),
			AllowedOriginPatterns: pointer.From(input.AllowedOriginPatterns),
			ExposedHeaders:        pointer.From(input.ExposedHeaders),
			MaxAgeSeconds:         pointer.From(input.MaxAge),
		},
	}
}

func flattenGatewayGatewayResourceRequests(input *appplatform.GatewayResourceRequests) []QuotaModel {
	if input == nil {
		return make([]QuotaModel, 0)
	}
	return []QuotaModel{
		{
			Cpu:    pointer.From(input.Cpu),
			Memory: pointer.From(input.Memory),
		},
	}
}

func flattenGatewaySsoProperties(input *appplatform.SsoProperties, old []GatewaySsoModel) []GatewaySsoModel {
	if input == nil {
		return make([]GatewaySsoModel, 0)
	}

	oldItems := make(map[string]GatewaySsoModel)
	for _, item := range old {
		oldItems[item.IssuerUri] = item
	}

	var issuerUri string
	if input.IssuerUri != nil {
		issuerUri = *input.IssuerUri
	}
	var clientId string
	var clientSecret string
	if oldItem, ok := oldItems[issuerUri]; ok {
		clientId = oldItem.ClientId
		clientSecret = oldItem.ClientSecret
	}
	return []GatewaySsoModel{
		{
			ClientId:     clientId,
			ClientSecret: clientSecret,
			IssuerUri:    issuerUri,
			Scope:        pointer.From(input.Scope),
		},
	}
}

func flattenGatewayGatewayApmTypes(input *[]appplatform.ApmType) []string {
	if input == nil {
		return nil
	}
	out := make([]string, 0)
	for _, v := range *input {
		out = append(out, string(v))
	}
	return out
}

func flattenGatewayClientAuth(input *appplatform.GatewayPropertiesClientAuth) []ClientAuthorizationModel {
	if input == nil || input.Certificates == nil || len(*input.Certificates) == 0 {
		return make([]ClientAuthorizationModel, 0)
	}
	certificateIds := make([]string, 0)
	if input.Certificates != nil {
		for _, v := range *input.Certificates {
			certId, err := appplatform.ParseCertificateIDInsensitively(v)
			if err == nil {
				certificateIds = append(certificateIds, certId.ID())
			}
		}
	}
	verificationEnabled := false
	if input.CertificateVerification != nil && *input.CertificateVerification == appplatform.GatewayCertificateVerificationEnabled {
		verificationEnabled = true
	}
	return []ClientAuthorizationModel{
		{
			CertificateIds:      certificateIds,
			VerificationEnabled: verificationEnabled,
		},
	}
}

func flattenGatewayLocalResponseCachePerRouteProperties(input appplatform.GatewayResponseCacheProperties) []ResponseCacheModel {
	if input == nil {
		return make([]ResponseCacheModel, 0)
	}
	if v, ok := input.(appplatform.GatewayLocalResponseCachePerRouteProperties); ok {
		return []ResponseCacheModel{
			{
				Size:       pointer.From(v.Size),
				TimeToLive: pointer.From(v.TimeToLive),
			},
		}
	}
	return make([]ResponseCacheModel, 0)
}

func flattenGatewayLocalResponseCachePerInstanceProperties(input appplatform.GatewayResponseCacheProperties) []ResponseCacheModel {
	if input == nil {
		return make([]ResponseCacheModel, 0)
	}
	if v, ok := input.(appplatform.GatewayLocalResponseCachePerInstanceProperties); ok {
		return []ResponseCacheModel{
			{
				Size:       pointer.From(v.Size),
				TimeToLive: pointer.From(v.TimeToLive),
			},
		}
	}
	return make([]ResponseCacheModel, 0)
}

func flattenGatewayApms(input *[]appplatform.ApmReference) ([]string, error) {
	out := make([]string, 0)
	if input == nil {
		return out, nil
	}
	for _, v := range *input {
		id, err := appplatform.ParseApmIDInsensitively(v.ResourceId)
		if err != nil {
			return nil, err
		}
		out = append(out, id.ID())
	}
	return out, nil
}
