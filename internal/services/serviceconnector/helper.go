// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serviceconnector

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AuthInfoModel struct {
	Type           string `tfschema:"type"`
	Name           string `tfschema:"name"`
	Secret         string `tfschema:"secret"`
	ClientId       string `tfschema:"client_id"`
	PrincipalId    string `tfschema:"principal_id"`
	SubscriptionId string `tfschema:"subscription_id"`
	Certificate    string `tfschema:"certificate"`
}

type SecretStoreModel struct {
	KeyVaultId string `tfschema:"key_vault_id"`
}

func secretStoreSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"key_vault_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func authInfoSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(servicelinker.AuthTypeSecret),
						string(servicelinker.AuthTypeServicePrincipalSecret),
						string(servicelinker.AuthTypeServicePrincipalCertificate),
						string(servicelinker.AuthTypeSystemAssignedIdentity),
						string(servicelinker.AuthTypeUserAssignedIdentity),
					}, false),
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"subscription_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"principal_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"certificate": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func expandServiceConnectorAuthInfoForCreate(input []AuthInfoModel) (servicelinker.AuthInfoBase, error) {
	if err := validateServiceConnectorAuthInfo(input); err != nil {
		return nil, err
	}

	if len(input) == 0 {
		return nil, nil
	}

	in := input[0]
	switch servicelinker.AuthType(in.Type) {
	case servicelinker.AuthTypeSecret:
		return servicelinker.SecretAuthInfo{
			Name: pointer.To(in.Name),
			SecretInfo: servicelinker.ValueSecretInfo{
				Value: pointer.To(in.Secret),
			},
		}, nil

	case servicelinker.AuthTypeServicePrincipalSecret:
		return servicelinker.ServicePrincipalSecretAuthInfo{
			ClientId:    in.ClientId,
			PrincipalId: in.PrincipalId,
			Secret:      in.Secret,
		}, nil

	case servicelinker.AuthTypeServicePrincipalCertificate:
		return servicelinker.ServicePrincipalCertificateAuthInfo{
			Certificate: in.Certificate,
			ClientId:    in.ClientId,
			PrincipalId: in.PrincipalId,
		}, nil

	case servicelinker.AuthTypeSystemAssignedIdentity:
		return servicelinker.SystemAssignedIdentityAuthInfo{}, nil

	case servicelinker.AuthTypeUserAssignedIdentity:
		return servicelinker.UserAssignedIdentityAuthInfo{
			ClientId:       pointer.To(in.ClientId),
			SubscriptionId: pointer.To(in.SubscriptionId),
		}, nil
	}

	return nil, fmt.Errorf("unrecognised authentication type: %q", in.Type)
}

func expandServiceConnectorAuthInfoForUpdate(input []AuthInfoModel) (links.AuthInfoBase, error) {
	if err := validateServiceConnectorAuthInfo(input); err != nil {
		return nil, err
	}

	if len(input) == 0 {
		return nil, nil
	}

	in := input[0]
	switch links.AuthType(in.Type) {
	case links.AuthTypeSecret:
		return links.SecretAuthInfo{
			Name: pointer.To(in.Name),
			SecretInfo: links.ValueSecretInfo{
				Value: pointer.To(in.Secret),
			},
		}, nil

	case links.AuthTypeServicePrincipalSecret:
		return links.ServicePrincipalSecretAuthInfo{
			ClientId:    in.ClientId,
			PrincipalId: in.PrincipalId,
			Secret:      in.Secret,
		}, nil

	case links.AuthTypeServicePrincipalCertificate:
		return links.ServicePrincipalCertificateAuthInfo{
			Certificate: in.Certificate,
			ClientId:    in.ClientId,
			PrincipalId: in.PrincipalId,
		}, nil

	case links.AuthTypeSystemAssignedIdentity:
		return links.SystemAssignedIdentityAuthInfo{}, nil

	case links.AuthTypeUserAssignedIdentity:
		return links.UserAssignedIdentityAuthInfo{
			ClientId:       pointer.To(in.ClientId),
			SubscriptionId: pointer.To(in.SubscriptionId),
		}, nil
	}

	return nil, fmt.Errorf("unrecognised authentication type: %q", in.Type)
}

func validateServiceConnectorAuthInfo(input []AuthInfoModel) error {
	if len(input) > 0 {
		authInfo := input[0]
		switch servicelinker.AuthType(authInfo.Type) {
		case servicelinker.AuthTypeSecret:
			if authInfo.ClientId != "" {
				return fmt.Errorf("`client_id` cannot be set when `type` is set to `Secret`")
			}
			if authInfo.SubscriptionId != "" {
				return fmt.Errorf("`subscription_id` cannot be set when `type` is set to `Secret`")
			}
			if authInfo.PrincipalId != "" {
				return fmt.Errorf("`principal_id` cannot be set when `type` is set to `Secret`")
			}
			if authInfo.Certificate != "" {
				return fmt.Errorf("`certificate` cannot be set when `type` is set to `Secret`")
			}
			if authInfo.Name != "" && authInfo.Secret == "" {
				return fmt.Errorf("`name` cannot be set when `secret` is empty")
			}
			if authInfo.Name == "" && authInfo.Secret != "" {
				return fmt.Errorf("`secret` cannot be set when `name` is empty")
			}

		case servicelinker.AuthTypeSystemAssignedIdentity:
			if authInfo.Name != "" || authInfo.Secret != "" || authInfo.ClientId != "" || authInfo.SubscriptionId != "" || authInfo.PrincipalId != "" || authInfo.Certificate != "" {
				return fmt.Errorf("no other authentication parameters should be set when `type` is set to `SystemIdentity`")
			}

		case servicelinker.AuthTypeServicePrincipalSecret:
			if authInfo.ClientId == "" {
				return fmt.Errorf("`client_id` must be specified when `type` is set to `ServicePrincipal`")
			}
			if authInfo.PrincipalId == "" {
				return fmt.Errorf("`principal_id` must be specified when `type` is set to `ServicePrincipal`")
			}
			if authInfo.Secret == "" {
				return fmt.Errorf("`secret` must be specified when `type` is set to `ServicePrincipal`")
			}
			if authInfo.SubscriptionId != "" {
				return fmt.Errorf("`subscription_id` cannot be set when `type` is set to `ServicePrincipal`")
			}
			if authInfo.Name != "" {
				return fmt.Errorf("`name` cannot be set when `type` is set to `ServicePrincipal`")
			}
			if authInfo.Certificate != "" {
				return fmt.Errorf("`certificate` cannot be set when `type` is set to `ServicePrincipal`")
			}

		case servicelinker.AuthTypeServicePrincipalCertificate:
			if authInfo.ClientId == "" {
				return fmt.Errorf("`client_id` must be specified when `type` is set to `ServicePrincipalCertificate`")
			}
			if authInfo.PrincipalId == "" {
				return fmt.Errorf("`principal_id` must be specified when `type` is set to `ServicePrincipalCertificate`")
			}
			if authInfo.Certificate == "" {
				return fmt.Errorf("`certificate` must be specified when `type` is set to `ServicePrincipalCertificate`")
			}
			if authInfo.SubscriptionId != "" {
				return fmt.Errorf("`subscription_id` cannot be set when `type` is set to `ServicePrincipalCertificate`")
			}
			if authInfo.Name != "" {
				return fmt.Errorf("`name` cannot be set when `type` is set to `ServicePrincipalCertificate`")
			}
			if authInfo.Secret != "" {
				return fmt.Errorf("`secret` cannot be set when `type` is set to `ServicePrincipalCertificate`")
			}

		case servicelinker.AuthTypeUserAssignedIdentity:
			if authInfo.PrincipalId != "" {
				return fmt.Errorf("`principal_id` cannot be set when `type` is set to `UserIdentity`")
			}
			if authInfo.Certificate != "" {
				return fmt.Errorf("`certificate` cannot be set when `type` is set to `UserIdentity`")
			}
			if authInfo.Name != "" {
				return fmt.Errorf("`name` cannot be set when `type` is set to `UserIdentity`")
			}
			if authInfo.Secret != "" {
				return fmt.Errorf("`secret` cannot be set when `type` is set to `UserIdentity`")
			}
			if authInfo.ClientId == "" && authInfo.SubscriptionId != "" {
				return fmt.Errorf("`subscription_id` cannot be set when `client_id` is empty")
			}
			if authInfo.ClientId != "" && authInfo.SubscriptionId == "" {
				return fmt.Errorf("`client_id` cannot be set when `subscription_id` is empty")
			}
		}
	}

	return nil
}

func expandSecretStore(input []SecretStoreModel) *servicelinker.SecretStore {
	if len(input) == 0 {
		return nil
	}
	v := input[0]

	keyVaultId := v.KeyVaultId
	return &servicelinker.SecretStore{
		KeyVaultId: utils.String(keyVaultId),
	}
}

func flattenServiceConnectorAuthInfo(input servicelinker.AuthInfoBase, pwd string) []AuthInfoModel {
	var authType string
	var name string
	var secret string
	var clientId string
	var principalId string
	var subscriptionId string
	var certificate string

	if value, ok := input.(servicelinker.SecretAuthInfo); ok {
		authType = string(servicelinker.AuthTypeSecret)
		if value.Name != nil {
			name = *value.Name
		}
		secret = pwd
	}

	if _, ok := input.(servicelinker.SystemAssignedIdentityAuthInfo); ok {
		authType = string(servicelinker.AuthTypeSystemAssignedIdentity)
	}

	if value, ok := input.(servicelinker.UserAssignedIdentityAuthInfo); ok {
		authType = string(servicelinker.AuthTypeUserAssignedIdentity)
		if value.ClientId != nil {
			clientId = *value.ClientId
		}
		if value.SubscriptionId != nil {
			subscriptionId = *value.SubscriptionId
		}
	}

	if value, ok := input.(servicelinker.ServicePrincipalSecretAuthInfo); ok {
		authType = string(servicelinker.AuthTypeServicePrincipalSecret)
		clientId = value.ClientId
		principalId = value.PrincipalId
		secret = pwd
	}

	if value, ok := input.(servicelinker.ServicePrincipalCertificateAuthInfo); ok {
		authType = string(servicelinker.AuthTypeServicePrincipalCertificate)
		certificate = value.Certificate
		clientId = value.ClientId
		principalId = value.PrincipalId
	}

	return []AuthInfoModel{
		{
			Type:           authType,
			Name:           name,
			Secret:         secret,
			ClientId:       clientId,
			PrincipalId:    principalId,
			SubscriptionId: subscriptionId,
			Certificate:    certificate,
		},
	}
}

// TODO: Only support Azure resource for now. Will include ConfluentBootstrapServer and ConfluentSchemaRegistry in the future.
func flattenTargetService(input servicelinker.TargetServiceBase) string {
	var targetServiceId string

	if value, ok := input.(servicelinker.AzureResource); ok {
		if value.Id != nil {
			targetServiceId = *value.Id
			if parsedId, err := parse.StorageAccountDefaultBlobID(targetServiceId); err == nil {
				storageAccountId := commonids.NewStorageAccountID(parsedId.SubscriptionId, parsedId.ResourceGroup, parsedId.StorageAccountName)
				targetServiceId = storageAccountId.ID()
			}
		}
	}

	return targetServiceId
}

func flattenSecretStore(input servicelinker.SecretStore) []SecretStoreModel {
	var keyVaultId string
	if input.KeyVaultId != nil {
		keyVaultId = *input.KeyVaultId
	}

	return []SecretStoreModel{
		{
			KeyVaultId: keyVaultId,
		},
	}
}
