// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cognitiveValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	managedHsmHelpers "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/helpers"
	managedHsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	managedHsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.ResourceWithUpdate = AIServices{}

var _ sdk.ResourceWithCustomImporter = AIServices{}

type AIServices struct{}

func (r AIServices) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := cognitiveservicesaccounts.ParseAccountID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.Cognitive.AccountsClient
		resp, err := client.AccountsGet(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Kind == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if !strings.EqualFold(*resp.Model.Kind, "AIServices") {
			return fmt.Errorf("importing %s: specified account is not of kind `AIServices`, got `%s`", id, *resp.Model.Kind)
		}

		return nil
	}
}

type VirtualNetworkRules struct {
	SubnetID                         string `tfschema:"subnet_id"`
	IgnoreMissingVnetServiceEndpoint bool   `tfschema:"ignore_missing_vnet_service_endpoint"`
}

type NetworkACLs struct {
	Bypass              string                `tfschema:"bypass"`
	DefaultAction       string                `tfschema:"default_action"`
	IpRules             []string              `tfschema:"ip_rules"`
	VirtualNetworkRules []VirtualNetworkRules `tfschema:"virtual_network_rules"`
}

type CustomerManagedKey struct {
	IdentityClientID string `tfschema:"identity_client_id"`
	KeyVaultKeyID    string `tfschema:"key_vault_key_id"`
	ManagedHsmKeyID  string `tfschema:"managed_hsm_key_id"`
}

type AIServicesModel struct {
	Name                            string                                     `tfschema:"name"`
	ResourceGroupName               string                                     `tfschema:"resource_group_name"`
	Location                        string                                     `tfschema:"location"`
	SkuName                         string                                     `tfschema:"sku_name"`
	CustomSubdomainName             string                                     `tfschema:"custom_subdomain_name"`
	CustomerManagedKey              []CustomerManagedKey                       `tfschema:"customer_managed_key"`
	Fqdns                           []string                                   `tfschema:"fqdns"`
	Identity                        []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	LocalAuthorizationEnabled       bool                                       `tfschema:"local_authentication_enabled"`
	NetworkACLs                     []NetworkACLs                              `tfschema:"network_acls"`
	OutboundNetworkAccessRestricted bool                                       `tfschema:"outbound_network_access_restricted"`
	PublicNetworkAccess             string                                     `tfschema:"public_network_access"`
	Tags                            map[string]string                          `tfschema:"tags"`
	Endpoint                        string                                     `tfschema:"endpoint"`
	PrimaryAccessKey                string                                     `tfschema:"primary_access_key"`
	SecondaryAccessKey              string                                     `tfschema:"secondary_access_key"`
}

func (AIServices) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveValidate.AccountName(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"F0", "F1", "S0", "S", "S1", "S2", "S3", "S4", "S5", "S6", "P0", "P1", "P2", "E0", "DC0",
			}, false),
		},

		"custom_subdomain_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						ExactlyOneOf: []string{"customer_managed_key.0.managed_hsm_key_id", "customer_managed_key.0.key_vault_key_id"},
					},

					"managed_hsm_key_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.Any(managedHsmValidate.ManagedHSMDataPlaneVersionedKeyID, managedHsmValidate.ManagedHSMDataPlaneVersionlessKeyID),
						ExactlyOneOf: []string{"customer_managed_key.0.managed_hsm_key_id", "customer_managed_key.0.key_vault_key_id"},
					},

					"identity_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"fqdns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"network_acls": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			RequiredWith: []string{"custom_subdomain_name"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bypass": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  cognitiveservicesaccounts.ByPassSelectionAzureServices,
						ValidateFunc: validation.StringInSlice(
							cognitiveservicesaccounts.PossibleValuesForByPassSelection(),
							false,
						),
					},
					"default_action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(cognitiveservicesaccounts.NetworkRuleActionAllow),
							string(cognitiveservicesaccounts.NetworkRuleActionDeny),
						}, false),
					},
					"ip_rules": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      set.HashIPv4AddressOrCIDR,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								commonValidate.IPv4Address,
								commonValidate.CIDR,
							),
						},
					},

					"virtual_network_rules": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"ignore_missing_vnet_service_endpoint": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
				},
			},
		},

		"outbound_network_access_restricted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(cognitiveservicesaccounts.PublicNetworkAccessEnabled),
				string(cognitiveservicesaccounts.PublicNetworkAccessDisabled),
			}, false),
			Default: string(cognitiveservicesaccounts.PublicNetworkAccessEnabled),
		},

		"storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},

					"identity_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (AIServices) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (AIServices) ModelObject() interface{} {
	return &AIServicesModel{}
}

func (AIServices) ResourceType() string {
	return "azurerm_ai_services"
}

func (AIServices) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AIServicesModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Cognitive.AccountsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := cognitiveservicesaccounts.NewAccountID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.AccountsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_ai_services", id.ID())
			}

			networkACLs, subnetIds := expandNetworkACLs(model.NetworkACLs)

			// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
			virtualNetworkNames := make([]string, 0)
			for _, v := range subnetIds {
				subnetId, err := commonids.ParseSubnetID(v)
				if err != nil {
					return err
				}
				if !utils.SliceContainsValue(virtualNetworkNames, subnetId.VirtualNetworkName) {
					virtualNetworkNames = append(virtualNetworkNames, subnetId.VirtualNetworkName)
				}
			}

			locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
			defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

			props := cognitiveservicesaccounts.Account{
				Kind:     pointer.To("AIServices"),
				Location: pointer.To(location.Normalize(model.Location)),
				Sku: &cognitiveservicesaccounts.Sku{
					Name: model.SkuName,
				},
				Properties: &cognitiveservicesaccounts.AccountProperties{
					NetworkAcls:                   networkACLs,
					CustomSubDomainName:           pointer.To(model.CustomSubdomainName),
					AllowedFqdnList:               pointer.To(model.Fqdns),
					PublicNetworkAccess:           pointer.To(cognitiveservicesaccounts.PublicNetworkAccess(model.PublicNetworkAccess)),
					RestrictOutboundNetworkAccess: pointer.To(model.OutboundNetworkAccessRestricted),
					DisableLocalAuth:              pointer.To(!model.LocalAuthorizationEnabled),
				},
				Tags: pointer.To(model.Tags),
			}

			expandIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			props.Identity = expandIdentity

			if err := client.AccountsCreateThenPoll(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// creating with KV HSM takes more time than expected, at least hours in most cases and eventually terminated by service
			if len(model.CustomerManagedKey) > 0 {
				customerManagedKey, err := expandCustomerManagedKey(model.CustomerManagedKey)
				if err != nil {
					return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
				}

				if customerManagedKey != nil {
					props.Properties.Encryption = customerManagedKey
					if err := client.AccountsUpdateThenPoll(ctx, id, props); err != nil {
						return fmt.Errorf("updating %s: %+v", id, err)
					}
				}
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (AIServices) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountsClient
			env := metadata.Client.Account.Environment

			state := AIServicesModel{}
			id, err := cognitiveservicesaccounts.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state.Name = id.AccountName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(model.Location)
				if sku := model.Sku; sku != nil {
					state.SkuName = sku.Name
				}

				identityFlatten, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = *identityFlatten

				if props := model.Properties; props != nil {
					state.Endpoint = pointer.From(props.Endpoint)
					state.CustomSubdomainName = pointer.From(props.CustomSubDomainName)
					state.NetworkACLs = flattenNetworkACLs(props.NetworkAcls)
					state.Fqdns = pointer.From(props.AllowedFqdnList)

					state.PublicNetworkAccess = string(pointer.From(props.PublicNetworkAccess))
					state.OutboundNetworkAccessRestricted = pointer.From(props.RestrictOutboundNetworkAccess)

					localAuthEnabled := true
					if props.DisableLocalAuth != nil {
						localAuthEnabled = !*props.DisableLocalAuth
					}

					if localAuthEnabled {
						keys, err := client.AccountsListKeys(ctx, *id)
						if err != nil {
							return fmt.Errorf("listing the Keys for %s: %+v", id, err)
						}

						if model := keys.Model; model != nil {
							state.PrimaryAccessKey = pointer.From(model.Key1)
							state.SecondaryAccessKey = pointer.From(model.Key2)
						}
					}
					state.LocalAuthorizationEnabled = localAuthEnabled

					customerManagedKey, err := flattenCustomerManagedKey(props.Encryption, env)
					if err != nil {
						return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
					}
					state.CustomerManagedKey = customerManagedKey
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (AIServices) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountsClient

			var model AIServicesModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := cognitiveservicesaccounts.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			props := resp.Model
			if metadata.ResourceData.HasChange("network_acls") {
				networkACLs, subnetIds := expandNetworkACLs(model.NetworkACLs)
				locks.MultipleByName(&subnetIds, network.VirtualNetworkResourceName)
				defer locks.UnlockMultipleByName(&subnetIds, network.VirtualNetworkResourceName)

				// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
				virtualNetworkNames := make([]string, 0)
				for _, v := range subnetIds {
					subnetId, err := commonids.ParseSubnetIDInsensitively(v)
					if err != nil {
						return err
					}
					if !utils.SliceContainsValue(virtualNetworkNames, subnetId.VirtualNetworkName) {
						virtualNetworkNames = append(virtualNetworkNames, subnetId.VirtualNetworkName)
					}
				}

				locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
				defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

				props.Properties.NetworkAcls = networkACLs
			}

			if metadata.ResourceData.HasChange("sku_name") {
				props.Sku = &cognitiveservicesaccounts.Sku{
					Name: model.SkuName,
				}
			}

			if metadata.ResourceData.HasChange("custom_subdomain_name") {
				props.Properties.CustomSubDomainName = pointer.To(model.CustomSubdomainName)
			}

			if metadata.ResourceData.HasChange("fqdns") {
				props.Properties.AllowedFqdnList = pointer.To(model.Fqdns)
			}

			if metadata.ResourceData.HasChange("public_network_access") {
				props.Properties.PublicNetworkAccess = pointer.To(cognitiveservicesaccounts.PublicNetworkAccess(model.PublicNetworkAccess))
			}

			if metadata.ResourceData.HasChange("outbound_network_access_restricted") {
				props.Properties.RestrictOutboundNetworkAccess = pointer.To(model.OutboundNetworkAccessRestricted)
			}

			if metadata.ResourceData.HasChange("local_authentication_enabled") {
				props.Properties.DisableLocalAuth = pointer.To(!model.LocalAuthorizationEnabled)
			}

			if metadata.ResourceData.HasChange("tags") {
				props.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("customer_managed_key") {
				customerManagedKey, err := expandCustomerManagedKey(model.CustomerManagedKey)
				if err != nil {
					return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
				}

				if customerManagedKey != nil {
					props.Properties.Encryption = customerManagedKey
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				expandIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				props.Identity = expandIdentity
			}

			if err := client.AccountsUpdateThenPoll(ctx, *id, *props); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (AIServices) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountsClient

			id, err := cognitiveservicesaccounts.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			log.Printf("[DEBUG] Retrieving %s..", *id)
			account, err := client.AccountsGet(ctx, *id)
			if err != nil || account.Model == nil || account.Model.Location == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			deletedAzureAIServicesId := cognitiveservicesaccounts.NewDeletedAccountID(id.SubscriptionId, *account.Model.Location, id.ResourceGroupName, id.AccountName)
			if err != nil {
				return err
			}

			log.Printf("[DEBUG] Deleting %s..", *id)
			if err := client.AccountsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			if metadata.Client.Features.CognitiveAccount.PurgeSoftDeleteOnDestroy {
				log.Printf("[DEBUG] Purging %s..", *id)
				if err := client.DeletedAccountsPurgeThenPoll(ctx, deletedAzureAIServicesId); err != nil {
					return fmt.Errorf("purging %s: %+v", *id, err)
				}
			} else {
				log.Printf("[DEBUG] Skipping Purge of %s", *id)
			}

			return nil
		},
	}
}

func (AIServices) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cognitiveservicesaccounts.ValidateAccountID
}

func expandCustomerManagedKey(input []CustomerManagedKey) (*cognitiveservicesaccounts.Encryption, error) {
	if len(input) == 0 {
		return &cognitiveservicesaccounts.Encryption{
			KeySource: pointer.To(cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices),
		}, nil
	}

	v := input[0]

	var identityClientId string
	if value := v.IdentityClientID; value != "" {
		identityClientId = value
	}

	encryption := &cognitiveservicesaccounts.Encryption{
		KeySource: pointer.To(cognitiveservicesaccounts.KeySourceMicrosoftPointKeyVault),
		KeyVaultProperties: &cognitiveservicesaccounts.KeyVaultProperties{
			IdentityClientId: pointer.To(identityClientId),
		},
	}

	if v.KeyVaultKeyID != "" {
		keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(v.KeyVaultKeyID)
		if err != nil {
			return nil, err
		}
		encryption.KeyVaultProperties.KeyName = pointer.To(keyId.Name)
		encryption.KeyVaultProperties.KeyVersion = pointer.To(keyId.Version)
		encryption.KeyVaultProperties.KeyVaultUri = pointer.To(keyId.KeyVaultBaseUrl)
	} else {
		hsmKyId, err := managedHsmParse.ManagedHSMDataPlaneVersionedKeyID(v.ManagedHsmKeyID, nil)
		if err != nil {
			return nil, err
		}

		encryption.KeyVaultProperties.KeyName = pointer.To(hsmKyId.KeyName)
		encryption.KeyVaultProperties.KeyVersion = pointer.To(hsmKyId.KeyVersion)
		encryption.KeyVaultProperties.KeyVaultUri = pointer.To(hsmKyId.BaseUri())
	}
	return encryption, nil
}

func flattenCustomerManagedKey(input *cognitiveservicesaccounts.Encryption, env environments.Environment) ([]CustomerManagedKey, error) {
	if input == nil || *input.KeySource == cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices {
		return []CustomerManagedKey{}, nil
	}

	keyName := ""
	keyVaultURI := ""
	keyVersion := ""
	customerManagerKey := CustomerManagedKey{}

	if props := input.KeyVaultProperties; props != nil {
		if props.KeyName != nil {
			keyName = *props.KeyName
		}
		if props.KeyVaultUri != nil {
			keyVaultURI = *props.KeyVaultUri
		}
		if props.KeyVersion != nil {
			keyVersion = *props.KeyVersion
		}

		isHsmURI, err, instanceName, domainSuffix := managedHsmHelpers.IsManagedHSMURI(env, keyVaultURI)
		if err != nil {
			return nil, err
		}

		if props.IdentityClientId != nil {
			customerManagerKey.IdentityClientID = *props.IdentityClientId
		}

		switch {
		case isHsmURI && keyVersion == "":
			{
				keyVaultKeyId := managedHsmParse.NewManagedHSMDataPlaneVersionlessKeyID(instanceName, domainSuffix, keyName)
				customerManagerKey.ManagedHsmKeyID = keyVaultKeyId.ID()
			}
		case isHsmURI && keyVersion != "":
			{
				keyVaultKeyId := managedHsmParse.NewManagedHSMDataPlaneVersionedKeyID(instanceName, domainSuffix, keyName, keyVersion)
				customerManagerKey.ManagedHsmKeyID = keyVaultKeyId.ID()
			}
		case !isHsmURI:
			{
				keyVaultKeyId, err := keyVaultParse.NewNestedItemID(keyVaultURI, keyVaultParse.NestedItemTypeKey, keyName, keyVersion)
				if err != nil {
					return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
				}
				customerManagerKey.KeyVaultKeyID = keyVaultKeyId.ID()
			}
		}
	}

	return []CustomerManagedKey{customerManagerKey}, nil
}

func expandNetworkACLs(input []NetworkACLs) (*cognitiveservicesaccounts.NetworkRuleSet, []string) {
	subnetIds := make([]string, 0)
	if len(input) == 0 {
		return nil, subnetIds
	}

	v := input[0]

	defaultAction := cognitiveservicesaccounts.NetworkRuleAction(v.DefaultAction)

	ipRules := make([]cognitiveservicesaccounts.IPRule, 0)

	for _, val := range v.IpRules {
		rule := cognitiveservicesaccounts.IPRule{
			Value: val,
		}
		ipRules = append(ipRules, rule)
	}

	networkRules := make([]cognitiveservicesaccounts.VirtualNetworkRule, 0)
	for _, val := range v.VirtualNetworkRules {
		subnetId := val.SubnetID
		subnetIds = append(subnetIds, subnetId)
		rule := cognitiveservicesaccounts.VirtualNetworkRule{
			Id:                               subnetId,
			IgnoreMissingVnetServiceEndpoint: pointer.To(val.IgnoreMissingVnetServiceEndpoint),
		}
		networkRules = append(networkRules, rule)
	}

	bypass := cognitiveservicesaccounts.ByPassSelection((v.Bypass))

	ruleSet := cognitiveservicesaccounts.NetworkRuleSet{
		Bypass:              &bypass,
		DefaultAction:       &defaultAction,
		IPRules:             &ipRules,
		VirtualNetworkRules: &networkRules,
	}
	return &ruleSet, subnetIds
}

func flattenNetworkACLs(input *cognitiveservicesaccounts.NetworkRuleSet) []NetworkACLs {
	if input == nil {
		return []NetworkACLs{}
	}

	ipRules := make([]string, 0)
	if input.IPRules != nil {
		for _, v := range *input.IPRules {
			ipRules = append(ipRules, v.Value)
		}
	}

	virtualNetworkRules := make([]VirtualNetworkRules, 0)
	if input.VirtualNetworkRules != nil {
		for _, v := range *input.VirtualNetworkRules {
			id := v.Id
			subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
			if err == nil {
				id = subnetId.ID()
			}

			virtualNetworkRules = append(virtualNetworkRules, VirtualNetworkRules{
				SubnetID:                         id,
				IgnoreMissingVnetServiceEndpoint: pointer.From(v.IgnoreMissingVnetServiceEndpoint),
			})
		}
	}

	return []NetworkACLs{{
		Bypass:              string(pointer.From(input.Bypass)),
		DefaultAction:       string(pointer.From(input.DefaultAction)),
		IpRules:             ipRules,
		VirtualNetworkRules: virtualNetworkRules,
	}}
}
