// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:hcp_openshift_cluster_id,resource_group_name:hcp_openshift_cluster_id,hcp_open_shift_cluster_name:hcp_openshift_cluster_id"

package redhatopenshift

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/externalauths"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithCustomizeDiff = RedHatHcpOpenShiftExternalAuth{}
	_ sdk.ResourceWithUpdate        = RedHatHcpOpenShiftExternalAuth{}
	_ sdk.ResourceWithIdentity      = RedHatHcpOpenShiftExternalAuth{}
)

type RedHatHcpOpenShiftExternalAuth struct{}

type RedHatHcpOpenShiftExternalAuthModel struct {
	Name      string                           `tfschema:"name"`
	ClusterId string                           `tfschema:"hcp_openshift_cluster_id"`
	Issuer    []HcpOpenShiftExternalAuthIssuer `tfschema:"issuer"`
	Claim     []HcpOpenShiftExternalAuthClaim  `tfschema:"claim"`
	Clients   []HcpOpenShiftExternalAuthClient `tfschema:"client"`
}

type HcpOpenShiftExternalAuthIssuer struct {
	Url       string   `tfschema:"url"`
	Audiences []string `tfschema:"audiences"`
	Ca        string   `tfschema:"ca"`
}

type HcpOpenShiftExternalAuthClaim struct {
	Mappings        []HcpOpenShiftExternalAuthMappings       `tfschema:"mappings"`
	ValidationRules []HcpOpenShiftExternalAuthValidationRule `tfschema:"validation_rule"`
}

type HcpOpenShiftExternalAuthMappings struct {
	Username []HcpOpenShiftExternalAuthUsername `tfschema:"username"`
	Groups   []HcpOpenShiftExternalAuthGroups   `tfschema:"groups"`
}

type HcpOpenShiftExternalAuthUsername struct {
	Claim        string `tfschema:"claim"`
	Prefix       string `tfschema:"prefix"`
	PrefixPolicy string `tfschema:"prefix_policy"`
}

type HcpOpenShiftExternalAuthGroups struct {
	Claim  string `tfschema:"claim"`
	Prefix string `tfschema:"prefix"`
}

type HcpOpenShiftExternalAuthValidationRule struct {
	Type          string                                  `tfschema:"type"`
	RequiredClaim []HcpOpenShiftExternalAuthRequiredClaim `tfschema:"required_claim"`
}

type HcpOpenShiftExternalAuthRequiredClaim struct {
	Claim         string `tfschema:"claim"`
	RequiredValue string `tfschema:"required_value"`
}

type HcpOpenShiftExternalAuthClient struct {
	ClientId    string                                    `tfschema:"client_id"`
	Type        string                                    `tfschema:"type"`
	ExtraScopes []string                                  `tfschema:"extra_scopes"`
	Component   []HcpOpenShiftExternalAuthClientComponent `tfschema:"component"`
}

type HcpOpenShiftExternalAuthClientComponent struct {
	Name                string `tfschema:"name"`
	AuthClientNamespace string `tfschema:"auth_client_namespace"`
}

func (RedHatHcpOpenShiftExternalAuth) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff
			const prefixPolicyPath = "claim.0.mappings.0.username.0.prefix_policy"
			const prefixPath = "claim.0.mappings.0.username.0.prefix"
			if diff.NewValueKnown(prefixPolicyPath) && diff.NewValueKnown(prefixPath) {
				if diff.Get(prefixPolicyPath).(string) == string(externalauths.UsernameClaimPrefixPolicyPrefix) && diff.Get(prefixPath).(string) == "" {
					return fmt.Errorf("`%s` is required when `%s` is `Prefix`", prefixPath, prefixPolicyPath)
				}
			}
			const rulesPath = "claim.0.validation_rule"
			if !diff.NewValueKnown(rulesPath + ".#") {
				return nil
			}
			for index := 0; index < diff.Get(rulesPath+".#").(int); index++ {
				typePath := fmt.Sprintf("%s.%d.type", rulesPath, index)
				claimPath := fmt.Sprintf("%s.%d.required_claim", rulesPath, index)
				if !diff.NewValueKnown(typePath) || !diff.NewValueKnown(claimPath+".#") {
					continue
				}
				if diff.Get(typePath).(string) == string(externalauths.TokenValidationRuleTypeRequiredClaim) && diff.Get(claimPath+".#").(int) == 0 {
					return fmt.Errorf("`%s` is required when `%s` is `RequiredClaim`", claimPath, typePath)
				}
			}
			return nil
		},
	}
}

func (r RedHatHcpOpenShiftExternalAuth) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hcp_openshift_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hcpopenshiftclusters.ValidateHcpOpenShiftClusterID,
		},

		"claim": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"mappings": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"username": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"claim": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"prefix": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"prefix_policy": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(externalauths.PossibleValuesForUsernameClaimPrefixPolicy(), false),
									Default:      externalauths.UsernameClaimPrefixPolicyNone,
								},
							}},
						},
						"groups": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"claim": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"prefix": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							}},
						},
					}},
				},
				"validation_rule": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(externalauths.PossibleValuesForTokenValidationRuleType(), false),
						},
						"required_claim": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"claim": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"required_value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							}},
						},
					}},
				},
			}},
		},

		"issuer": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},
				"audiences": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					MaxItems: 10,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"ca": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			}},
		},

		"client": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 20,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"client_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(externalauths.PossibleValuesForExternalAuthClientType(), false),
				},
				"extra_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"component": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"auth_client_namespace": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					}},
				},
			}},
		},
	}
}

func (r RedHatHcpOpenShiftExternalAuth) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RedHatHcpOpenShiftExternalAuth) ModelObject() any {
	return &RedHatHcpOpenShiftExternalAuthModel{}
}

func (r RedHatHcpOpenShiftExternalAuth) ResourceType() string {
	return "azurerm_redhat_hcp_openshift_external_auth"
}

func (r RedHatHcpOpenShiftExternalAuth) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return externalauths.ValidateExternalAuthID
}

func (r RedHatHcpOpenShiftExternalAuth) Identity() resourceids.ResourceId {
	return &externalauths.ExternalAuthId{}
}

func (r RedHatHcpOpenShiftExternalAuth) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config RedHatHcpOpenShiftExternalAuthModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			parent, err := hcpopenshiftclusters.ParseHcpOpenShiftClusterID(config.ClusterId)
			if err != nil {
				return err
			}
			id := externalauths.NewExternalAuthID(parent.SubscriptionId, parent.ResourceGroupName, parent.HcpOpenShiftClusterName, config.Name)
			client := metadata.Client.RedHatOpenShift.ExternalAuthsClient
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}
			parameters, err := expandHcpOpenShiftExternalAuth(config)
			if err != nil {
				return err
			}
			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r RedHatHcpOpenShiftExternalAuth) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := externalauths.ParseExternalAuthID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			var config RedHatHcpOpenShiftExternalAuthModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			parameters, err := expandHcpOpenShiftExternalAuth(config)
			if err != nil {
				return err
			}
			if err := metadata.Client.RedHatOpenShift.ExternalAuthsClient.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r RedHatHcpOpenShiftExternalAuth) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := externalauths.ParseExternalAuthID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := metadata.Client.RedHatOpenShift.ExternalAuthsClient.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			return r.flatten(metadata, *id, *resp.Model)
		},
	}
}

func (r RedHatHcpOpenShiftExternalAuth) flatten(metadata sdk.ResourceMetaData, id externalauths.ExternalAuthId, model externalauths.ExternalAuth) error {
	state, err := flattenHcpOpenShiftExternalAuth(id, model)
	if err != nil {
		return err
	}
	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
		return err
	}
	return metadata.Encode(state)
}

func (r RedHatHcpOpenShiftExternalAuth) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := externalauths.ParseExternalAuthID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := metadata.Client.RedHatOpenShift.ExternalAuthsClient.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}

func expandHcpOpenShiftExternalAuth(input RedHatHcpOpenShiftExternalAuthModel) (*externalauths.ExternalAuth, error) {
	if len(input.Issuer) != 1 || len(input.Claim) != 1 || len(input.Claim[0].Mappings) != 1 || len(input.Claim[0].Mappings[0].Username) != 1 {
		return nil, fmt.Errorf("expanding external auth: exactly one `issuer`, `claim`, `claim.mappings` and `claim.mappings.username` block is required")
	}
	issuer := input.Issuer[0]
	mappings := input.Claim[0].Mappings[0]
	username := mappings.Username[0]
	if username.PrefixPolicy == string(externalauths.UsernameClaimPrefixPolicyPrefix) && username.Prefix == "" {
		return nil, fmt.Errorf("`claim.0.mappings.0.username.0.prefix` is required when `claim.0.mappings.0.username.0.prefix_policy` is `Prefix`")
	}
	properties := &externalauths.ExternalAuthProperties{
		Issuer: externalauths.TokenIssuerProfile{
			Url:       issuer.Url,
			Audiences: issuer.Audiences,
		},
		Claim: externalauths.ExternalAuthClaimProfile{
			Mappings: externalauths.TokenClaimMappingsProfile{
				Username: externalauths.UsernameClaimProfile{
					Claim: username.Claim,
				},
			},
		},
	}
	if issuer.Ca != "" {
		properties.Issuer.Ca = pointer.To(issuer.Ca)
	}
	if username.Prefix != "" {
		properties.Claim.Mappings.Username.Prefix = pointer.To(username.Prefix)
	}
	if username.PrefixPolicy != "" {
		properties.Claim.Mappings.Username.PrefixPolicy = pointer.ToEnum[externalauths.UsernameClaimPrefixPolicy](username.PrefixPolicy)
	}
	if len(mappings.Groups) > 0 {
		groups := mappings.Groups[0]
		properties.Claim.Mappings.Groups = &externalauths.GroupClaimProfile{
			Claim: groups.Claim,
		}
		if groups.Prefix != "" {
			properties.Claim.Mappings.Groups.Prefix = pointer.To(groups.Prefix)
		}
	}
	rules := make([]externalauths.TokenClaimValidationRule, 0, len(input.Claim[0].ValidationRules))
	for _, rule := range input.Claim[0].ValidationRules {
		if rule.Type == string(externalauths.TokenValidationRuleTypeRequiredClaim) && len(rule.RequiredClaim) != 1 {
			return nil, fmt.Errorf("expanding `claim.validation_rule.required_claim`: exactly one block is required when `type` is `RequiredClaim`")
		}
		validationRule := externalauths.TokenClaimValidationRule{
			Type: pointer.ToEnum[externalauths.TokenValidationRuleType](rule.Type),
		}
		if len(rule.RequiredClaim) > 0 {
			validationRule.RequiredClaim = &externalauths.TokenRequiredClaim{
				Claim:         rule.RequiredClaim[0].Claim,
				RequiredValue: rule.RequiredClaim[0].RequiredValue,
			}
		}
		rules = append(rules, validationRule)
	}
	properties.Claim.ValidationRules = &rules
	clients := make([]externalauths.ExternalAuthClientProfile, 0, len(input.Clients))
	for _, client := range input.Clients {
		if len(client.Component) != 1 {
			return nil, fmt.Errorf("expanding `client.component`: exactly one block is required")
		}
		scopes := append([]string{}, client.ExtraScopes...)
		clients = append(clients, externalauths.ExternalAuthClientProfile{
			ClientId:    client.ClientId,
			Type:        externalauths.ExternalAuthClientType(client.Type),
			ExtraScopes: &scopes,
			Component: externalauths.ExternalAuthClientComponentProfile{
				Name:                client.Component[0].Name,
				AuthClientNamespace: client.Component[0].AuthClientNamespace,
			},
		})
	}
	properties.Clients = &clients
	return &externalauths.ExternalAuth{
		Properties: properties,
	}, nil
}

func flattenHcpOpenShiftExternalAuth(id externalauths.ExternalAuthId, input externalauths.ExternalAuth) (*RedHatHcpOpenShiftExternalAuthModel, error) {
	if input.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	parent := hcpopenshiftclusters.NewHcpOpenShiftClusterID(id.SubscriptionId, id.ResourceGroupName, id.HcpOpenShiftClusterName)
	properties := input.Properties
	state := &RedHatHcpOpenShiftExternalAuthModel{
		Name:      id.ExternalAuthName,
		ClusterId: parent.ID(),
		Issuer: []HcpOpenShiftExternalAuthIssuer{
			{
				Url:       properties.Issuer.Url,
				Audiences: properties.Issuer.Audiences,
				Ca:        pointer.From(properties.Issuer.Ca),
			},
		},
	}
	username := properties.Claim.Mappings.Username
	mappings := HcpOpenShiftExternalAuthMappings{
		Username: []HcpOpenShiftExternalAuthUsername{
			{
				Claim:        username.Claim,
				Prefix:       pointer.From(username.Prefix),
				PrefixPolicy: pointer.FromEnum(username.PrefixPolicy),
			},
		},
	}
	if groups := properties.Claim.Mappings.Groups; groups != nil {
		mappings.Groups = []HcpOpenShiftExternalAuthGroups{
			{
				Claim:  groups.Claim,
				Prefix: pointer.From(groups.Prefix),
			},
		}
	}
	claim := HcpOpenShiftExternalAuthClaim{
		Mappings: []HcpOpenShiftExternalAuthMappings{mappings},
	}
	for _, rule := range pointer.From(properties.Claim.ValidationRules) {
		validationRule := HcpOpenShiftExternalAuthValidationRule{
			Type: pointer.FromEnum(rule.Type),
		}
		if rule.RequiredClaim != nil {
			validationRule.RequiredClaim = []HcpOpenShiftExternalAuthRequiredClaim{
				{
					Claim:         rule.RequiredClaim.Claim,
					RequiredValue: rule.RequiredClaim.RequiredValue,
				},
			}
		}
		claim.ValidationRules = append(claim.ValidationRules, validationRule)
	}
	state.Claim = []HcpOpenShiftExternalAuthClaim{claim}
	for _, client := range pointer.From(properties.Clients) {
		state.Clients = append(state.Clients, HcpOpenShiftExternalAuthClient{
			ClientId:    client.ClientId,
			Type:        string(client.Type),
			ExtraScopes: pointer.From(client.ExtraScopes),
			Component: []HcpOpenShiftExternalAuthClientComponent{
				{
					Name:                client.Component.Name,
					AuthClientNamespace: client.Component.AuthClientNamespace,
				},
			},
		})
	}
	return state, nil
}
