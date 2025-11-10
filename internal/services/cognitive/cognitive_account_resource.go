// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	search "github.com/hashicorp/go-azure-sdk/resource-manager/search/2025-05-01/services"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCognitiveAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCognitiveAccountCreate,
		Read:   resourceCognitiveAccountRead,
		Update: resourceCognitiveAccountUpdate,
		Delete: resourceCognitiveAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cognitiveservicesaccounts.ParseAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AIServices",
					"Academic",
					"AnomalyDetector",
					"Bing.Autosuggest",
					"Bing.Autosuggest.v7",
					"Bing.CustomSearch",
					"Bing.Search",
					"Bing.Search.v7",
					"Bing.Speech",
					"Bing.SpellCheck",
					"Bing.SpellCheck.v7",
					"CognitiveServices",
					"ComputerVision",
					"ContentModerator",
					"ConversationalLanguageUnderstanding",
					"ContentSafety",
					"CustomSpeech",
					"CustomVision.Prediction",
					"CustomVision.Training",
					"Emotion",
					"Face",
					"FormRecognizer",
					"ImmersiveReader",
					"LUIS",
					"LUIS.Authoring",
					"MetricsAdvisor",
					"OpenAI",
					"Personalizer",
					"QnAMaker",
					"Recommendations",
					"SpeakerRecognition",
					"Speech",
					"SpeechServices",
					"SpeechTranslation",
					"TextAnalytics",
					"TextTranslation",
					"WebLM",
				}, false),
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"C2", "C3", "C4", "D3", "DC0", "E0", "F0", "F1", "P0", "P1", "P2", "S", "S0", "S1", "S2", "S3", "S4", "S5", "S6",
				}, false),
			},

			"project_management_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"custom_subdomain_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
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
							Required:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
						},

						"identity_client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"dynamic_throttling_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"metrics_advisor_aad_client_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"metrics_advisor_aad_tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"metrics_advisor_super_user_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"metrics_advisor_website_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"network_acls": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"custom_subdomain_name"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
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
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.Any(
									commonValidate.IPv4Address,
									commonValidate.CIDR,
								),
							},
							Set: set.HashIPv4AddressOrCIDR,
						},

						"virtual_network_rules": {
							Type:       pluginsdk.TypeSet,
							Optional:   true,
							ConfigMode: pluginsdk.SchemaConfigModeAuto,
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

						"bypass": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								cognitiveservicesaccounts.PossibleValuesForByPassSelection(),
								false,
							),
						},
					},
				},
			},

			"network_injection": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"scenario": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cognitiveservicesaccounts.ScenarioTypeAgent),
							}, false),
						},

						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: commonids.ValidateSubnetID,
						},
					},
				},
			},

			"outbound_network_access_restricted": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"qna_runtime_endpoint": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"custom_question_answering_search_service_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: search.ValidateSearchServiceID,
			},

			"custom_question_answering_search_service_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Sensitive:    true,
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
		},

		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i interface{}) error {
			kind := d.Get("kind").(string)

			if d.Get("project_management_enabled").(bool) {
				if kind != "AIServices" {
					return errors.New("`project_management_enabled` can only be set to `true` when `kind` is set to `AIServices`")
				}

				if len(d.Get("identity").([]interface{})) == 0 {
					return errors.New("for `project_management_enabled` to be set to `true`, a managed identity must be assigned. Please configure the `identity` block")
				}

				if d.HasChange("customer_managed_key") && len(d.Get("customer_managed_key").([]interface{})) == 0 {
					if err := d.ForceNew("customer_managed_key"); err != nil {
						return err
					}
				}
			} else if d.HasChange("project_management_enabled") && kind != "OpenAI" {
				if err := d.ForceNew("project_management_enabled"); err != nil {
					return err
				}
			}

			if d.Get("dynamic_throttling_enabled").(bool) && slices.Contains([]string{"OpenAI", "AIServices"}, kind) {
				return errors.New("`dynamic_throttling_enabled` is currently not supported when `kind` is set to `OpenAI` or `AIServices`")
			}

			if bypass, ok := d.GetOk("network_acls.0.bypass"); ok && bypass != "" && !slices.Contains([]string{"OpenAI", "AIServices", "TextAnalytics"}, kind) {
				return fmt.Errorf("`network_acls.bypass` cannot be set when `kind` is set to `%s`", kind)
			}

			networkInjection := d.Get("network_injection").([]interface{})
			if len(networkInjection) > 0 && networkInjection[0] != nil {
				if kind != "AIServices" {
					return errors.New("the `network_injection` block is only supported when `kind` is set to `AIServices`")
				}
			}

			if d.HasChange("custom_subdomain_name") {
				old, _ := d.GetChange("custom_subdomain_name")
				if old != nil && old != "" {
					if err := d.ForceNew("custom_subdomain_name"); err != nil {
						return err
					}
				}
			}

			if d.HasChange("kind") {
				old, new := d.GetChange("kind")
				oldKind := old.(string)
				newKind := new.(string)

				// Only allow changing `kind` from/to `OpenAI` or `AIServices`, force new for all others
				if !slices.Contains([]string{"OpenAI", "AIServices"}, newKind) || !slices.Contains([]string{"OpenAI", "AIServices"}, oldKind) {
					if err := d.ForceNew("kind"); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

func resourceCognitiveAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cognitiveservicesaccounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.AccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cognitive_account", id.ID())
		}
	}

	sku := cognitiveservicesaccounts.Sku{
		Name: d.Get("sku_name").(string),
	}

	networkAcls, subnetIds := expandCognitiveAccountNetworkAcls(d)

	// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err := commonids.ParseSubnetIDInsensitively(v)
		if err != nil {
			return err
		}
		if !slices.Contains(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	publicNetworkAccess := cognitiveservicesaccounts.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = cognitiveservicesaccounts.PublicNetworkAccessDisabled
	}

	apiProps, err := expandCognitiveAccountAPIProperties(d)
	if err != nil {
		return err
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props := cognitiveservicesaccounts.Account{
		Kind:     pointer.To(d.Get("kind").(string)),
		Location: pointer.To(azure.NormalizeLocation(d.Get("location").(string))),
		Sku:      &sku,
		Properties: &cognitiveservicesaccounts.AccountProperties{
			ApiProperties:                 apiProps,
			NetworkAcls:                   networkAcls,
			CustomSubDomainName:           pointer.To(d.Get("custom_subdomain_name").(string)),
			AllowedFqdnList:               utils.ExpandStringSlice(d.Get("fqdns").([]interface{})),
			PublicNetworkAccess:           &publicNetworkAccess,
			UserOwnedStorage:              expandCognitiveAccountStorage(d.Get("storage").([]interface{})),
			RestrictOutboundNetworkAccess: pointer.To(d.Get("outbound_network_access_restricted").(bool)),
			DisableLocalAuth:              pointer.To(!d.Get("local_auth_enabled").(bool)),
			DynamicThrottlingEnabled:      pointer.To(d.Get("dynamic_throttling_enabled").(bool)),
			AllowProjectManagement:        pointer.To(d.Get("project_management_enabled").(bool)),
			NetworkInjections:             expandCognitiveAccountNetworkInjection(d.Get("network_injection").([]interface{})),
			Encryption:                    expandCognitiveAccountCustomerManagedKey(d.Get("customer_managed_key").([]interface{})),
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.AccountsCreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCognitiveAccountRead(d, meta)
}

func resourceCognitiveAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	sku := cognitiveservicesaccounts.Sku{
		Name: d.Get("sku_name").(string),
	}

	networkAcls, subnetIds := expandCognitiveAccountNetworkAcls(d)

	// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err := commonids.ParseSubnetIDInsensitively(v)
		if err != nil {
			return err
		}
		if !slices.Contains(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	publicNetworkAccess := cognitiveservicesaccounts.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = cognitiveservicesaccounts.PublicNetworkAccessDisabled
	}

	apiProps, err := expandCognitiveAccountAPIProperties(d)
	if err != nil {
		return err
	}

	props := cognitiveservicesaccounts.Account{
		Sku:  &sku,
		Kind: pointer.To(d.Get("kind").(string)),
		Properties: &cognitiveservicesaccounts.AccountProperties{
			ApiProperties:                 apiProps,
			NetworkAcls:                   networkAcls,
			CustomSubDomainName:           pointer.To(d.Get("custom_subdomain_name").(string)),
			AllowedFqdnList:               utils.ExpandStringSlice(d.Get("fqdns").([]interface{})),
			PublicNetworkAccess:           &publicNetworkAccess,
			UserOwnedStorage:              expandCognitiveAccountStorage(d.Get("storage").([]interface{})),
			RestrictOutboundNetworkAccess: pointer.To(d.Get("outbound_network_access_restricted").(bool)),
			DisableLocalAuth:              pointer.To(!d.Get("local_auth_enabled").(bool)),
			DynamicThrottlingEnabled:      pointer.To(d.Get("dynamic_throttling_enabled").(bool)),
			AllowProjectManagement:        pointer.To(d.Get("project_management_enabled").(bool)),
			NetworkInjections:             expandCognitiveAccountNetworkInjection(d.Get("network_injection").([]interface{})),
			Encryption:                    expandCognitiveAccountCustomerManagedKey(d.Get("customer_managed_key").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if d.HasChanges("customer_managed_key") {
		old, new := d.GetChange("customer_managed_key")
		// Remove `customer_managed_key` (switch using a customer managed key to microsoft managed), and explicitly specify KeySource as `Microsoft.CognitiveServices`.
		if len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
			props.Properties.Encryption = &cognitiveservicesaccounts.Encryption{
				KeySource: pointer.To(cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices),
			}
		}
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		props.Identity = identity
	}

	if err = client.AccountsUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCognitiveAccountRead(d, meta)
}

func resourceCognitiveAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("kind", model.Kind)

		d.Set("location", location.NormalizeNilable(model.Location))
		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		d.Set("identity", identity)

		if props := model.Properties; props != nil {
			if apiProps := props.ApiProperties; apiProps != nil {
				d.Set("qna_runtime_endpoint", apiProps.QnaRuntimeEndpoint)
				d.Set("custom_question_answering_search_service_id", apiProps.QnaAzureSearchEndpointId)
				d.Set("metrics_advisor_aad_client_id", apiProps.AadClientId)
				d.Set("metrics_advisor_aad_tenant_id", apiProps.AadTenantId)
				d.Set("metrics_advisor_super_user_name", apiProps.SuperUser)
				d.Set("metrics_advisor_website_name", apiProps.WebsiteName)
			}
			d.Set("endpoint", props.Endpoint)
			d.Set("custom_subdomain_name", props.CustomSubDomainName)
			if err := d.Set("network_acls", flattenCognitiveAccountNetworkAcls(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_acls` for Cognitive Account %q: %+v", id, err)
			}

			networkInjection, err := flattenCognitiveAccountNetworkInjection(props.NetworkInjections)
			if err != nil {
				return err
			}

			if err := d.Set("network_injection", networkInjection); err != nil {
				return fmt.Errorf("setting `network_injection`: %+v", err)
			}

			dynamicThrottlingEnabled := false
			if props.DynamicThrottlingEnabled != nil {
				dynamicThrottlingEnabled = *props.DynamicThrottlingEnabled
			}
			d.Set("dynamic_throttling_enabled", dynamicThrottlingEnabled)

			d.Set("fqdns", pointer.From(props.AllowedFqdnList))

			// Azure API issue: `AllowProjectManagement` not reset during the rollback, see: https://github.com/Azure/azure-rest-api-specs/issues/38678
			allowProjectManagement := pointer.From(props.AllowProjectManagement)
			if pointer.From(model.Kind) == "OpenAI" {
				allowProjectManagement = false
			}
			d.Set("project_management_enabled", allowProjectManagement)

			publicNetworkAccess := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = *props.PublicNetworkAccess == cognitiveservicesaccounts.PublicNetworkAccessEnabled
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)

			if err := d.Set("storage", flattenCognitiveAccountStorage(props.UserOwnedStorage)); err != nil {
				return fmt.Errorf("setting `storages` for Cognitive Account %q: %+v", id, err)
			}
			outboundNetworkAccessRestricted := false
			if props.RestrictOutboundNetworkAccess != nil {
				outboundNetworkAccessRestricted = *props.RestrictOutboundNetworkAccess
			}
			// lintignore:R001
			d.Set("outbound_network_access_restricted", outboundNetworkAccessRestricted)

			localAuthEnabled := true
			if props.DisableLocalAuth != nil {
				localAuthEnabled = !*props.DisableLocalAuth
			}
			d.Set("local_auth_enabled", localAuthEnabled)

			if localAuthEnabled {
				keys, err := client.AccountsListKeys(ctx, *id)
				if err != nil {
					// note for the resource we shouldn't gracefully fail since we have permission to CRUD it
					return fmt.Errorf("listing the Keys for %s: %+v", *id, err)
				}

				if model := keys.Model; model != nil {
					d.Set("primary_access_key", model.Key1)
					d.Set("secondary_access_key", model.Key2)
				}
			}

			customerManagedKey, err := flattenCognitiveAccountCustomerManagedKey(props.Encryption)
			if err != nil {
				return err
			}

			if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceCognitiveAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Cognitive.AccountsClient
	deletedAccountsClient := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	// first we need to retrieve it, since we need the location to be able to purge it
	log.Printf("[DEBUG] Retrieving %s..", *id)
	account, err := accountsClient.AccountsGet(ctx, *id)
	if err != nil || account.Model == nil || account.Model.Location == nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	deletedAccountId := cognitiveservicesaccounts.NewDeletedAccountID(id.SubscriptionId, *account.Model.Location, id.ResourceGroupName, id.AccountName)

	log.Printf("[DEBUG] Deleting %s..", *id)
	if err := accountsClient.AccountsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if meta.(*clients.Client).Features.CognitiveAccount.PurgeSoftDeleteOnDestroy {
		log.Printf("[DEBUG] Purging %s..", *id)
		if err := deletedAccountsClient.DeletedAccountsPurgeThenPoll(ctx, deletedAccountId); err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}
	} else {
		log.Printf("[DEBUG] Skipping Purge of %s", *id)
	}

	// If `network_injection` is configured, wait for Service Association Link (SAL) to be removed from the agent subnet
	// This is a service issue workaround as the SAL is not removed immediately after the Cognitive Account is deleted
	var subnetId *commonids.SubnetId
	if account.Model.Properties != nil && account.Model.Properties.NetworkInjections != nil {
		networkInjections := *account.Model.Properties.NetworkInjections
		if len(networkInjections) > 0 {
			if subnetIdStr := pointer.From(networkInjections[0].SubnetArmId); subnetIdStr != "" {
				parsedSubnetId, err := commonids.ParseSubnetIDInsensitively(subnetIdStr)
				if err != nil {
					return err
				}
				subnetId = parsedSubnetId
			}
		}
	}

	if subnetId != nil {
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("internal-error: context had no deadline")
		}

		subnetClient := meta.(*clients.Client).Network.Subnets
		log.Printf("[DEBUG] Waiting for Service Association Links to be removed from subnet %s", subnetId.ID())

		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{"SALExists"},
			Target:     []string{"SALDeleted"},
			Refresh:    serviceAssociationLinkStateRefreshFunc(ctx, subnetClient, *subnetId),
			MinTimeout: 15 * time.Second,
			Timeout:    time.Until(deadline),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for Service Association Links to be removed from %s: %+v", *subnetId, err)
		}
	}

	return nil
}

func serviceAssociationLinkStateRefreshFunc(ctx context.Context, client *subnets.SubnetsClient, subnetId commonids.SubnetId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, subnetId, subnets.DefaultGetOperationOptions())
		if err != nil {
			return nil, "Error", fmt.Errorf("retrieving subnet %s: %+v", subnetId, err)
		}

		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.ServiceAssociationLinks != nil {
			serviceAssociationLinks := *resp.Model.Properties.ServiceAssociationLinks

			if len(serviceAssociationLinks) > 0 {
				for _, sal := range serviceAssociationLinks {
					// The SAL name associated with the Cognitive Services account is "legionservicelink"
					if pointer.From(sal.Name) == "legionservicelink" && sal.Properties != nil && sal.Properties.ProvisioningState != nil {
						log.Printf("[DEBUG] Found Service Association Link %s with provisioning state %s",
							pointer.From(sal.Name), *sal.Properties.ProvisioningState)
						return resp, "SALExists", nil
					}
				}
			}
		}

		log.Printf("[DEBUG] No Service Association Links found on subnet %s", subnetId.ID())
		return resp, "SALDeleted", nil
	}
}

func cognitiveAccountStateRefreshFunc(ctx context.Context, client *cognitiveservicesaccounts.CognitiveServicesAccountsClient, id cognitiveservicesaccounts.AccountId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.AccountsGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties != nil && res.Model.Properties.ProvisioningState != nil {
			return res, string(*res.Model.Properties.ProvisioningState), nil
		}
		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}

func expandCognitiveAccountNetworkAcls(d *pluginsdk.ResourceData) (*cognitiveservicesaccounts.NetworkRuleSet, []string) {
	input := d.Get("network_acls").([]interface{})
	subnetIds := make([]string, 0)
	if len(input) == 0 || input[0] == nil {
		return nil, subnetIds
	}

	v := input[0].(map[string]interface{})

	defaultAction := cognitiveservicesaccounts.NetworkRuleAction(v["default_action"].(string))

	ipRulesRaw := v["ip_rules"].(*pluginsdk.Set)
	ipRules := make([]cognitiveservicesaccounts.IPRule, 0)

	for _, v := range ipRulesRaw.List() {
		rule := cognitiveservicesaccounts.IPRule{
			Value: v.(string),
		}
		ipRules = append(ipRules, rule)
	}

	networkRules := make([]cognitiveservicesaccounts.VirtualNetworkRule, 0)
	networkRulesRaw := v["virtual_network_rules"]
	for _, v := range networkRulesRaw.(*pluginsdk.Set).List() {
		value := v.(map[string]interface{})
		subnetId := value["subnet_id"].(string)
		subnetIds = append(subnetIds, subnetId)
		rule := cognitiveservicesaccounts.VirtualNetworkRule{
			Id:                               subnetId,
			IgnoreMissingVnetServiceEndpoint: pointer.To(value["ignore_missing_vnet_service_endpoint"].(bool)),
		}
		networkRules = append(networkRules, rule)
	}

	ruleSet := cognitiveservicesaccounts.NetworkRuleSet{
		DefaultAction:       &defaultAction,
		IPRules:             &ipRules,
		VirtualNetworkRules: &networkRules,
	}

	if b, ok := d.GetOk("network_acls.0.bypass"); ok && b != "" {
		bypasss := cognitiveservicesaccounts.ByPassSelection(v["bypass"].(string))
		ruleSet.Bypass = &bypasss
	}

	return &ruleSet, subnetIds
}

func expandCognitiveAccountStorage(input []interface{}) *[]cognitiveservicesaccounts.UserOwnedStorage {
	if len(input) == 0 {
		return nil
	}
	results := make([]cognitiveservicesaccounts.UserOwnedStorage, 0)
	for _, v := range input {
		value := v.(map[string]interface{})
		results = append(results, cognitiveservicesaccounts.UserOwnedStorage{
			ResourceId:       pointer.To(value["storage_account_id"].(string)),
			IdentityClientId: pointer.To(value["identity_client_id"].(string)),
		})
	}
	return &results
}

func expandCognitiveAccountAPIProperties(d *pluginsdk.ResourceData) (*cognitiveservicesaccounts.ApiProperties, error) {
	props := cognitiveservicesaccounts.ApiProperties{}
	kind := d.Get("kind")
	if kind == "QnAMaker" {
		if v, ok := d.GetOk("qna_runtime_endpoint"); ok && v != "" {
			props.QnaRuntimeEndpoint = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("the QnAMaker runtime endpoint `qna_runtime_endpoint` is required when kind is set to `QnAMaker`")
		}
	}
	if v, ok := d.GetOk("custom_question_answering_search_service_id"); ok {
		if kind == "TextAnalytics" {
			props.QnaAzureSearchEndpointId = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("the Search Service ID `custom_question_answering_search_service_id` can only be set when kind is set to `TextAnalytics`")
		}
	}
	if v, ok := d.GetOk("custom_question_answering_search_service_key"); ok {
		if kind == "TextAnalytics" {
			props.QnaAzureSearchEndpointKey = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("the Search Service Key `custom_question_answering_search_service_key` can only be set when kind is set to `TextAnalytics`")
		}
	}

	if v, ok := d.GetOk("metrics_advisor_aad_client_id"); ok {
		if kind == "MetricsAdvisor" {
			props.AadClientId = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_aad_client_id can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_aad_tenant_id"); ok {
		if kind == "MetricsAdvisor" {
			props.AadTenantId = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_aad_tenant_id can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_super_user_name"); ok {
		if kind == "MetricsAdvisor" {
			props.SuperUser = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_super_user_name can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	if v, ok := d.GetOk("metrics_advisor_website_name"); ok {
		if kind == "MetricsAdvisor" {
			props.WebsiteName = pointer.To(v.(string))
		} else {
			return nil, fmt.Errorf("metrics_advisor_website_name can only used set when kind is set to `MetricsAdvisor`")
		}
	}
	return &props, nil
}

func flattenCognitiveAccountNetworkAcls(input *cognitiveservicesaccounts.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ipRules := make([]interface{}, 0)
	if input.IPRules != nil {
		for _, v := range *input.IPRules {
			ipRules = append(ipRules, v.Value)
		}
	}

	virtualNetworkRules := make([]interface{}, 0)
	if input.VirtualNetworkRules != nil {
		for _, v := range *input.VirtualNetworkRules {
			id := v.Id
			subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
			if err == nil {
				id = subnetId.ID()
			}

			virtualNetworkRules = append(virtualNetworkRules, map[string]interface{}{
				"subnet_id":                            id,
				"ignore_missing_vnet_service_endpoint": *v.IgnoreMissingVnetServiceEndpoint,
			})
		}
	}

	return []interface{}{map[string]interface{}{
		"bypass":                input.Bypass,
		"default_action":        input.DefaultAction,
		"ip_rules":              pluginsdk.NewSet(pluginsdk.HashString, ipRules),
		"virtual_network_rules": virtualNetworkRules,
	}}
}

func flattenCognitiveAccountStorage(input *[]cognitiveservicesaccounts.UserOwnedStorage) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	results := make([]interface{}, 0)
	for _, v := range *input {
		value := make(map[string]interface{})
		if v.ResourceId != nil {
			value["storage_account_id"] = *v.ResourceId
		}
		if v.IdentityClientId != nil {
			value["identity_client_id"] = *v.IdentityClientId
		}
		results = append(results, value)
	}
	return results
}

func expandCognitiveAccountCustomerManagedKey(input []interface{}) *cognitiveservicesaccounts.Encryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	keyId, _ := keyVaultParse.ParseOptionallyVersionedNestedItemID(v["key_vault_key_id"].(string))
	if keyId == nil {
		return nil
	}

	keySource := cognitiveservicesaccounts.KeySourceMicrosoftPointKeyVault

	var identity string
	if value := v["identity_client_id"]; value != nil && value != "" {
		identity = value.(string)
	}

	return &cognitiveservicesaccounts.Encryption{
		KeySource: &keySource,
		KeyVaultProperties: &cognitiveservicesaccounts.KeyVaultProperties{
			KeyName:          pointer.To(keyId.Name),
			KeyVersion:       pointer.To(keyId.Version),
			KeyVaultUri:      pointer.To(keyId.KeyVaultBaseUrl),
			IdentityClientId: pointer.To(identity),
		},
	}
}

func flattenCognitiveAccountCustomerManagedKey(input *cognitiveservicesaccounts.Encryption) ([]interface{}, error) {
	if input == nil || pointer.From(input.KeySource) == cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices {
		return []interface{}{}, nil
	}

	var keyId string
	var identityClientId string
	if props := input.KeyVaultProperties; props != nil {
		keyVaultKeyId, err := keyVaultParse.NewNestedItemID(*props.KeyVaultUri, keyVaultParse.NestedItemTypeKey, *props.KeyName, *props.KeyVersion)
		if err != nil {
			return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
		}
		keyId = keyVaultKeyId.ID()
		if props.IdentityClientId != nil {
			identityClientId = *props.IdentityClientId
		}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_key_id":   keyId,
			"identity_client_id": identityClientId,
		},
	}, nil
}

func expandCognitiveAccountNetworkInjection(input []interface{}) *[]cognitiveservicesaccounts.NetworkInjection {
	if len(input) == 0 {
		return nil
	}

	results := make([]cognitiveservicesaccounts.NetworkInjection, 0)
	for _, v := range input {
		m := v.(map[string]interface{})

		scenario := cognitiveservicesaccounts.ScenarioType(m["scenario"].(string))

		var subnetId *string
		if m["subnet_id"] != nil && m["subnet_id"] != "" {
			subnetId = pointer.To(m["subnet_id"].(string))
		}

		results = append(results, cognitiveservicesaccounts.NetworkInjection{
			Scenario:    &scenario,
			SubnetArmId: subnetId,
		})
	}

	return &results
}

func flattenCognitiveAccountNetworkInjection(input *[]cognitiveservicesaccounts.NetworkInjection) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		var subnetId string
		if v.SubnetArmId != nil {
			subnet, err := commonids.ParseSubnetIDInsensitively(*v.SubnetArmId)
			if err != nil {
				return nil, err
			}
			subnetId = subnet.ID()
		}

		results = append(results, map[string]interface{}{
			"scenario":  v.Scenario,
			"subnet_id": subnetId,
		})
	}

	return results, nil
}
