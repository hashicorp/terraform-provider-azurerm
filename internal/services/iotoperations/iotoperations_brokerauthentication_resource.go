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
	Name                   string                              `tfschema:"name"`
	ResourceGroupName      string                              `tfschema:"resource_group_name"`
	InstanceName           string                              `tfschema:"instance_name"`
	BrokerName             string                              `tfschema:"broker_name"`
	AuthenticationMethods  []BrokerAuthenticationMethodModel   `tfschema:"authentication_methods"`
	ExtendedLocation       *ExtendedLocationModel              `tfschema:"extended_location"`
	Tags                   map[string]string                   `tfschema:"tags"`
	ProvisioningState      *string                             `tfschema:"provisioning_state"`
}

type BrokerAuthenticationMethodModel struct {
	Method         string                                    `tfschema:"method"`
	CustomSettings *BrokerAuthenticationCustomSettingsModel `tfschema:"custom_settings"`
}

type BrokerAuthenticationCustomSettingsModel struct {
	Auth      *BrokerAuthenticationAuthModel `tfschema:"auth"`
	CaCertPem *string                        `tfschema:"ca_cert_pem"`
}

type BrokerAuthenticationAuthModel struct {
	X509 *BrokerAuthenticationX509Model `tfschema:"x509"`
}

type BrokerAuthenticationX509Model struct {
	SecretRef           *string                                           `tfschema:"secret_ref"`
	KeyVault            *BrokerAuthenticationKeyVaultModel                `tfschema:"key_vault"`
	SecretProviderClass *BrokerAuthenticationSecretProviderClassModel     `tfschema:"secret_provider_class"`
	Attributes          map[string]string                                 `tfschema:"attributes"`
}

type BrokerAuthenticationKeyVaultModel struct {
	Vault                  *BrokerAuthenticationKeyVaultVaultModel          `tfschema:"vault"`
	SecretProviderClass    *BrokerAuthenticationSecretProviderClassModel     `tfschema:"secret_provider_class"`
}

type BrokerAuthenticationKeyVaultVaultModel struct {
	Credentials   *BrokerAuthenticationKeyVaultCredentialsModel `tfschema:"credentials"`
	DirectoryId   *string                                       `tfschema:"directory_id"`
	Name          *string                                       `tfschema:"name"`
}

type BrokerAuthenticationKeyVaultCredentialsModel struct {
	ServicePrincipalLocalSecretRef *string `tfschema:"service_principal_local_secret_ref"`
}

type BrokerAuthenticationSecretProviderClassModel struct {
	Spec *BrokerAuthenticationSecretProviderClassSpecModel `tfschema:"spec"`
}

type BrokerAuthenticationSecretProviderClassSpecModel struct {
	AzureKeyVault *BrokerAuthenticationAzureKeyVaultModel `tfschema:"azure_key_vault"`
}

type BrokerAuthenticationAzureKeyVaultModel struct {
	KeyvaultName *string                                             `tfschema:"keyvault_name"`
	Secrets      []BrokerAuthenticationAzureKeyVaultSecretModel     `tfschema:"secrets"`
	TenantId     *string                                             `tfschema:"tenant_id"`
}

type BrokerAuthenticationAzureKeyVaultSecretModel struct {
	Name     *string `tfschema:"name"`
	SecretId *string `tfschema:"secret_id"`
	Version  *string `tfschema:"version"`
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
		"authentication_methods": {
			Type:     pluginsdk.TypeList,
			Required: true,
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
								"auth": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"x509": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"secret_ref": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
														"attributes": {
															Type:     pluginsdk.TypeMap,
															Optional: true,
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
								"ca_cert_pem": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"CustomLocation",
						}, false),
					},
				},
			},
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r BrokerAuthenticationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			// NOTE: O+C Azure automatically assigns provisioning state during resource lifecycle
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

			// Build payload
			payload := brokerauthentication.BrokerAuthenticationResource{
				Properties: expandBrokerAuthenticationProperties(model.AuthenticationMethods),
			}

			if model.ExtendedLocation != nil {
				payload.ExtendedLocation = expandExtendedLocation(model.ExtendedLocation)
			}

			if len(model.Tags) > 0 {
				payload.Tags = &model.Tags
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
				if respModel.ExtendedLocation != nil {
					model.ExtendedLocation = flattenExtendedLocation(respModel.ExtendedLocation)
				}

				if respModel.Tags != nil {
					model.Tags = *respModel.Tags
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

			// Check if anything actually changed before making API call
			if !metadata.ResourceData.HasChange("tags") && 
			   !metadata.ResourceData.HasChange("authentication_methods") {
				return nil
			}

			payload := brokerauthentication.BrokerAuthenticationPatchModel{}
			hasChanges := false

			// Only include tags if they changed
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			// Only include authentication methods if they changed
			if metadata.ResourceData.HasChange("authentication_methods") {
				payload.Properties = &brokerauthentication.BrokerAuthenticatorPropertiesPatch{
					AuthenticationMethods: expandBrokerAuthenticationMethods(model.AuthenticationMethods),
				}
				hasChanges = true
			}

			// Only make API call if something actually changed
			if !hasChanges {
				return nil
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
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
func expandBrokerAuthenticationProperties(methods []BrokerAuthenticationMethodModel) *brokerauthentication.BrokerAuthenticatorProperties {
	if len(methods) == 0 {
		return &brokerauthentication.BrokerAuthenticatorProperties{}
	}

	result := &brokerauthentication.BrokerAuthenticatorProperties{
		AuthenticationMethods: expandBrokerAuthenticationMethods(methods),
	}

	return result
}

func expandBrokerAuthenticationMethods(methods []BrokerAuthenticationMethodModel) *[]brokerauthentication.BrokerAuthenticatorMethod {
	if len(methods) == 0 {
		return nil
	}

	result := make([]brokerauthentication.BrokerAuthenticatorMethod, 0, len(methods))

	for _, method := range methods {
		authMethod := brokerauthentication.BrokerAuthenticatorMethod{
			Method: method.Method,
		}

		if method.CustomSettings != nil {
			// Add custom settings expansion logic here when needed
		}

		result = append(result, authMethod)
	}

	return &result
}

func flattenBrokerAuthenticationProperties(props *brokerauthentication.BrokerAuthenticatorProperties) []BrokerAuthenticationMethodModel {
	if props == nil || props.AuthenticationMethods == nil {
		return []BrokerAuthenticationMethodModel{}
	}

	result := make([]BrokerAuthenticationMethodModel, 0, len(*props.AuthenticationMethods))

	for _, method := range *props.AuthenticationMethods {
		authMethod := BrokerAuthenticationMethodModel{
			Method: method.Method,
		}

		// Add custom settings flattening logic here when needed

		result = append(result, authMethod)
	}

	return result
}