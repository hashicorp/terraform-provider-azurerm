// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serviceconnector

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/servicelinker"
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

func expandServiceConnectorAuthInfo(input []AuthInfoModel) (servicelinker.AuthInfoBase, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("authentication should be defined")
	}
	v := input[0]

	authType := servicelinker.AuthType(v.Type)
	name := v.Name
	secret := v.Secret
	clientId := v.ClientId
	subscriptionId := v.SubscriptionId
	principalId := v.PrincipalId
	certificate := v.Certificate

	switch authType {
	case servicelinker.AuthTypeSecret:
		if clientId != "" {
			return nil, fmt.Errorf("`client_id` cannot be set when `type` is set to `Secret`")
		}
		if subscriptionId != "" {
			return nil, fmt.Errorf("`subscription_id` cannot be set when `type` is set to `Secret`")
		}
		if principalId != "" {
			return nil, fmt.Errorf("`principal_id` cannot be set when `type` is set to `Secret`")
		}
		if certificate != "" {
			return nil, fmt.Errorf("`certificate` cannot be set when `type` is set to `Secret`")
		}
		if name != "" && secret == "" {
			return nil, fmt.Errorf("`name` cannot be set when `secret` is empty")
		}
		if name == "" && secret != "" {
			return nil, fmt.Errorf("`secret` cannot be set when `name` is empty")
		}
		return servicelinker.SecretAuthInfo{
			Name: utils.String(name),
			SecretInfo: servicelinker.ValueSecretInfo{
				Value: utils.String(secret),
			},
		}, nil

	case servicelinker.AuthTypeSystemAssignedIdentity:
		if name != "" || secret != "" || clientId != "" || subscriptionId != "" || principalId != "" || certificate != "" {
			return nil, fmt.Errorf("no other parameters should be set when `type` is set to `SystemIdentity`")
		}
		return servicelinker.SystemAssignedIdentityAuthInfo{}, nil

	case servicelinker.AuthTypeServicePrincipalSecret:
		if clientId == "" {
			return nil, fmt.Errorf("`client_id` must be specified when `type` is set to `ServicePrincipal`")
		}
		if principalId == "" {
			return nil, fmt.Errorf("`principal_id` must be specified when `type` is set to `ServicePrincipal`")
		}
		if secret == "" {
			return nil, fmt.Errorf("`secret` must be specified when `type` is set to `ServicePrincipal`")
		}
		if subscriptionId != "" {
			return nil, fmt.Errorf("`subscription_id` cannot be set when `type` is set to `ServicePrincipal`")
		}
		if name != "" {
			return nil, fmt.Errorf("`name` cannot be set when `type` is set to `ServicePrincipal`")
		}
		if certificate != "" {
			return nil, fmt.Errorf("`certificate` cannot be set when `type` is set to `ServicePrincipal`")
		}
		return servicelinker.ServicePrincipalSecretAuthInfo{
			ClientId:    clientId,
			PrincipalId: principalId,
			Secret:      secret,
		}, nil

	case servicelinker.AuthTypeServicePrincipalCertificate:
		if clientId == "" {
			return nil, fmt.Errorf("`client_id` must be specified when `type` is set to `ServicePrincipalCertificate`")
		}
		if principalId == "" {
			return nil, fmt.Errorf("`principal_id` must be specified when `type` is set to `ServicePrincipalCertificate`")
		}
		if certificate == "" {
			return nil, fmt.Errorf("`certificate` must be specified when `type` is set to `ServicePrincipalCertificate`")
		}
		if subscriptionId != "" {
			return nil, fmt.Errorf("`subscription_id` cannot be set when `type` is set to `ServicePrincipalCertificate`")
		}
		if name != "" {
			return nil, fmt.Errorf("`name` cannot be set when `type` is set to `ServicePrincipalCertificate`")
		}
		if secret != "" {
			return nil, fmt.Errorf("`secret` cannot be set when `type` is set to `ServicePrincipalCertificate`")
		}
		return servicelinker.ServicePrincipalCertificateAuthInfo{
			Certificate: certificate,
			ClientId:    clientId,
			PrincipalId: principalId,
		}, nil

	case servicelinker.AuthTypeUserAssignedIdentity:
		if principalId != "" {
			return nil, fmt.Errorf("`principal_id` cannot be set when `type` is set to `UserIdentity`")
		}
		if certificate != "" {
			return nil, fmt.Errorf("`certificate` cannot be set when `type` is set to `UserIdentity`")
		}
		if name != "" {
			return nil, fmt.Errorf("`name` cannot be set when `type` is set to `UserIdentity`")
		}
		if secret != "" {
			return nil, fmt.Errorf("`secret` cannot be set when `type` is set to `UserIdentity`")
		}
		if clientId == "" && subscriptionId != "" {
			return nil, fmt.Errorf("`subscription_id` cannot be set when `client_id` is empty")
		}
		if clientId != "" && subscriptionId == "" {
			return nil, fmt.Errorf("`client_id` cannot be set when `subscription_id` is empty")
		}
		return servicelinker.UserAssignedIdentityAuthInfo{
			ClientId:       utils.String(clientId),
			SubscriptionId: utils.String(subscriptionId),
		}, nil
	}

	return nil, fmt.Errorf("unsupported authentication type %q", authType)
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
