package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BrokerAuthenticationResource struct{}

var _ sdk.ResourceWithUpdate = BrokerAuthenticationResource{}

type BrokerAuthenticationModel struct {
	Name                  string                            `tfschema:"name"`
	ResourceGroupName     string                            `tfschema:"resource_group_name"`
	InstanceName          string                            `tfschema:"instance_name"`
	BrokerName            string                            `tfschema:"broker_name"`
	ExtendedLocation      ExtendedLocationModel             `tfschema:"extended_location"`
	AuthenticationMethods []BrokerAuthenticationMethodModel `tfschema:"authentication_methods"`
	ProvisioningState     *string                           `tfschema:"provisioning_state"`
}

type BrokerAuthenticationMethodModel struct {
	Method                      string                                        `tfschema:"method"`
	CustomSettings              *BrokerAuthenticationCustomSettingsModel      `tfschema:"custom_settings"`
	ServiceAccountTokenSettings *BrokerAuthenticationServiceAccountTokenModel `tfschema:"service_account_token_settings"`
	X509Settings                *BrokerAuthenticationX509SettingsModel        `tfschema:"x509_settings"`
}

type BrokerAuthenticationCustomSettingsModel struct {
	Auth            *BrokerAuthenticationCustomAuthModel `tfschema:"auth"`
	CaCertConfigMap *string                              `tfschema:"ca_cert_config_map"`
	Endpoint        string                               `tfschema:"endpoint"`
	Headers         map[string]string                    `tfschema:"headers"`
}

type BrokerAuthenticationCustomAuthModel struct {
	X509 BrokerAuthenticationX509ManualModel `tfschema:"x509"`
}

type BrokerAuthenticationX509ManualModel struct {
	SecretRef string `tfschema:"secret_ref"`
}

type BrokerAuthenticationServiceAccountTokenModel struct {
	Audiences []string `tfschema:"audiences"`
}

type BrokerAuthenticationX509SettingsModel struct {
	AuthorizationAttributes map[string]BrokerAuthenticationX509AttributesModel `tfschema:"authorization_attributes"`
	TrustedClientCaCert     *string                                            `tfschema:"trusted_client_ca_cert"`
}

type BrokerAuthenticationX509AttributesModel struct {
	Attributes map[string]string `tfschema:"attributes"`
	Subject    string            `tfschema:"subject"`
}

func (r BrokerAuthenticationResource) ModelObject() interface{} {
	return &BrokerAuthenticationModel{}
}

func (r BrokerAuthenticationResource) ResourceType() string {
	return "azurerm_iotoperations_broker_authentication"
}

func (r BrokerAuthenticationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return brokerauthentication.ValidateAuthenticationID
}

func (r BrokerAuthenticationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 90),
		},
		"instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"broker_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Required: true, // Changed from optional since SDK requires it
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						Default:  "CustomLocation",
						ValidateFunc: validation.StringInSlice([]string{
							"CustomLocation",
						}, false),
					},
				},
			},
		},
		"authentication_methods": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"method": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Custom",
							"ServiceAccountToken",
							"X509",
						}, false),
					},
					"custom_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"auth": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"x509": {
												Type:     pluginsdk.TypeList,
												Required: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"secret_ref": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},
										},
									},
								},
								"ca_cert_config_map": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"headers": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"service_account_token_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"audiences": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"x509_settings": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"trusted_client_ca_cert": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"authorization_attributes": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"subject": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"attributes": {
												Type:     pluginsdk.TypeMap,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r BrokerAuthenticationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r BrokerAuthenticationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthenticationClient

			var model BrokerAuthenticationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := brokerauthentication.NewAuthenticationID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.BrokerName, model.Name)

			// Check if resource already exists
			existing, err := client.Get(ctx, id)
			if err == nil && existing.Model != nil {
				return fmt.Errorf("IoT Operations Broker Authentication %q already exists", id.AuthenticationName)
			}

			// Build payload with proper ExtendedLocation struct
			payload := brokerauthentication.BrokerAuthenticationResource{
				ExtendedLocation: brokerauthentication.ExtendedLocation{
					Name: *model.ExtendedLocation.Name,
					Type: brokerauthentication.ExtendedLocationType(*model.ExtendedLocation.Type),
				},
				Properties: expandBrokerAuthenticationProperties(model.AuthenticationMethods),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BrokerAuthenticationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthenticationClient

			id, err := brokerauthentication.ParseAuthenticationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := BrokerAuthenticationModel{
				Name:              id.AuthenticationName,
				ResourceGroupName: id.ResourceGroupName,
				InstanceName:      id.InstanceName,
				BrokerName:        id.BrokerName,
			}

			if respModel := resp.Model; respModel != nil {
				// Properly map ExtendedLocation struct
				model.ExtendedLocation = ExtendedLocationModel{
					Name: &respModel.ExtendedLocation.Name,
					Type: func() *string {
						s := string(respModel.ExtendedLocation.Type)
						return &s
					}(),
				}

				if respModel.Properties != nil {
					model.AuthenticationMethods = flattenBrokerAuthenticationProperties(respModel.Properties)

					if respModel.Properties.ProvisioningState != nil {
						provisioningState := string(*respModel.Properties.ProvisioningState)
						model.ProvisioningState = &provisioningState
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r BrokerAuthenticationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthenticationClient

			id, err := brokerauthentication.ParseAuthenticationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BrokerAuthenticationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Since there's no separate Update method, use CreateOrUpdate
			payload := brokerauthentication.BrokerAuthenticationResource{
				ExtendedLocation: brokerauthentication.ExtendedLocation{
					Name: *model.ExtendedLocation.Name,
					Type: brokerauthentication.ExtendedLocationType(*model.ExtendedLocation.Type),
				},
				Properties: expandBrokerAuthenticationProperties(model.AuthenticationMethods),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BrokerAuthenticationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerAuthenticationClient

			id, err := brokerauthentication.ParseAuthenticationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

// Helper functions for expand/flatten operations
func expandBrokerAuthenticationProperties(methods []BrokerAuthenticationMethodModel) *brokerauthentication.BrokerAuthenticationProperties {
	return &brokerauthentication.BrokerAuthenticationProperties{
		AuthenticationMethods: expandBrokerAuthenticationMethods(methods),
	}
}

func expandBrokerAuthenticationMethods(methods []BrokerAuthenticationMethodModel) []brokerauthentication.BrokerAuthenticatorMethods {
	result := make([]brokerauthentication.BrokerAuthenticatorMethods, 0, len(methods))

	for _, method := range methods {
		authMethod := brokerauthentication.BrokerAuthenticatorMethods{
			Method: brokerauthentication.BrokerAuthenticationMethod(method.Method),
		}

		if method.CustomSettings != nil {
			authMethod.CustomSettings = expandBrokerAuthenticationCustomSettings(*method.CustomSettings)
		}

		if method.ServiceAccountTokenSettings != nil {
			authMethod.ServiceAccountTokenSettings = expandBrokerAuthenticationServiceAccountToken(*method.ServiceAccountTokenSettings)
		}

		if method.X509Settings != nil {
			authMethod.X509Settings = expandBrokerAuthenticationX509Settings(*method.X509Settings)
		}

		result = append(result, authMethod)
	}

	return result
}

func expandBrokerAuthenticationCustomSettings(settings BrokerAuthenticationCustomSettingsModel) *brokerauthentication.BrokerAuthenticatorMethodCustom {
	result := &brokerauthentication.BrokerAuthenticatorMethodCustom{
		Endpoint: settings.Endpoint,
	}

	if settings.Auth != nil {
		result.Auth = &brokerauthentication.BrokerAuthenticatorCustomAuth{
			X509: brokerauthentication.X509ManualCertificate{
				SecretRef: settings.Auth.X509.SecretRef,
			},
		}
	}

	if settings.CaCertConfigMap != nil {
		result.CaCertConfigMap = settings.CaCertConfigMap
	}

	if len(settings.Headers) > 0 {
		result.Headers = &settings.Headers
	}

	return result
}

func expandBrokerAuthenticationServiceAccountToken(settings BrokerAuthenticationServiceAccountTokenModel) *brokerauthentication.BrokerAuthenticatorMethodSat {
	return &brokerauthentication.BrokerAuthenticatorMethodSat{
		Audiences: settings.Audiences,
	}
}

func expandBrokerAuthenticationX509Settings(settings BrokerAuthenticationX509SettingsModel) *brokerauthentication.BrokerAuthenticatorMethodX509 {
	result := &brokerauthentication.BrokerAuthenticatorMethodX509{}

	if settings.TrustedClientCaCert != nil {
		result.TrustedClientCaCert = settings.TrustedClientCaCert
	}

	if len(settings.AuthorizationAttributes) > 0 {
		authzAttrs := make(map[string]brokerauthentication.BrokerAuthenticatorMethodX509Attributes)
		for key, attr := range settings.AuthorizationAttributes {
			authzAttrs[key] = brokerauthentication.BrokerAuthenticatorMethodX509Attributes{
				Subject:    attr.Subject,
				Attributes: attr.Attributes,
			}
		}
		result.AuthorizationAttributes = &authzAttrs
	}

	return result
}

func flattenBrokerAuthenticationProperties(props *brokerauthentication.BrokerAuthenticationProperties) []BrokerAuthenticationMethodModel {
	if props == nil {
		return []BrokerAuthenticationMethodModel{}
	}

	result := make([]BrokerAuthenticationMethodModel, 0, len(props.AuthenticationMethods))

	for _, method := range props.AuthenticationMethods {
		authMethod := BrokerAuthenticationMethodModel{
			Method: string(method.Method),
		}

		if method.CustomSettings != nil {
			authMethod.CustomSettings = flattenBrokerAuthenticationCustomSettings(method.CustomSettings)
		}

		if method.ServiceAccountTokenSettings != nil {
			authMethod.ServiceAccountTokenSettings = flattenBrokerAuthenticationServiceAccountToken(method.ServiceAccountTokenSettings)
		}

		if method.X509Settings != nil {
			authMethod.X509Settings = flattenBrokerAuthenticationX509Settings(method.X509Settings)
		}

		result = append(result, authMethod)
	}

	return result
}

func flattenBrokerAuthenticationCustomSettings(settings *brokerauthentication.BrokerAuthenticatorMethodCustom) *BrokerAuthenticationCustomSettingsModel {
	result := &BrokerAuthenticationCustomSettingsModel{
		Endpoint: settings.Endpoint,
	}

	if settings.Auth != nil {
		result.Auth = &BrokerAuthenticationCustomAuthModel{
			X509: BrokerAuthenticationX509ManualModel{
				SecretRef: settings.Auth.X509.SecretRef,
			},
		}
	}

	if settings.CaCertConfigMap != nil {
		result.CaCertConfigMap = settings.CaCertConfigMap
	}

	if settings.Headers != nil {
		result.Headers = *settings.Headers
	}

	return result
}

func flattenBrokerAuthenticationServiceAccountToken(settings *brokerauthentication.BrokerAuthenticatorMethodSat) *BrokerAuthenticationServiceAccountTokenModel {
	return &BrokerAuthenticationServiceAccountTokenModel{
		Audiences: settings.Audiences,
	}
}

func flattenBrokerAuthenticationX509Settings(settings *brokerauthentication.BrokerAuthenticatorMethodX509) *BrokerAuthenticationX509SettingsModel {
	result := &BrokerAuthenticationX509SettingsModel{
		AuthorizationAttributes: make(map[string]BrokerAuthenticationX509AttributesModel),
	}

	if settings.TrustedClientCaCert != nil {
		result.TrustedClientCaCert = settings.TrustedClientCaCert
	}

	if settings.AuthorizationAttributes != nil {
		for key, attr := range *settings.AuthorizationAttributes {
			result.AuthorizationAttributes[key] = BrokerAuthenticationX509AttributesModel{
				Subject:    attr.Subject,
				Attributes: attr.Attributes,
			}
		}
	}

	return result
}
