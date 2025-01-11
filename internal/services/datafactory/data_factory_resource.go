// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedvirtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactory() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryCreateUpdate,
		Read:   resourceDataFactoryRead,
		Update: resourceDataFactoryCreateUpdate,
		Delete: resourceDataFactoryDelete,

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataFactoryV0ToV1{},
			1: migration.DataFactoryV1ToV2{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := factories.ParseFactoryID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"github_configuration": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"vsts_configuration"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"branch_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"git_url": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"repository_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"root_folder": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"publishing_enabled": {
							Type:     pluginsdk.TypeBool,
							Default:  true,
							Optional: true,
						},
					},
				},
			},

			"vsts_configuration": {
				Type:          pluginsdk.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"github_configuration"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"branch_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"project_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"repository_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"root_folder": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"publishing_enabled": {
							Type:     pluginsdk.TypeBool,
							Default:  true,
							Optional: true,
						},
					},
				},
			},

			"global_parameter": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(factories.PossibleValuesForGlobalParameterType(), false),
						},

						"value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"managed_virtual_network_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"public_network_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"purview_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: account.ValidateAccountID,
			},

			"customer_managed_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.NestedItemId,
			},

			"customer_managed_key_identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
				RequiredWith: []string{"customer_managed_key_id"},
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("managed_virtual_network_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(bool) && !new.(bool)
			}),
		),
	}
}

func resourceDataFactoryCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.Factories
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := factories.NewFactoryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, factories.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_factory", id.ID())
		}
	}

	publicNetworkAccess := factories.PublicNetworkAccessEnabled
	enabled := d.Get("public_network_enabled").(bool)
	if !enabled {
		publicNetworkAccess = factories.PublicNetworkAccessDisabled
	}

	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	location := location.Normalize(d.Get("location").(string))
	payload := factories.Factory{
		Location: utils.String(location),
		Properties: &factories.FactoryProperties{
			PublicNetworkAccess: &publicNetworkAccess,
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if purviewId, ok := d.GetOk("purview_id"); ok {
		payload.Properties.PurviewConfiguration = &factories.PurviewConfiguration{
			PurviewResourceId: pointer.To(purviewId.(string)),
		}
	}

	if keyVaultKeyID, ok := d.GetOk("customer_managed_key_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyID.(string))
		if err != nil {
			return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}

		payload.Properties.Encryption = &factories.EncryptionConfiguration{
			VaultBaseURL: keyVaultKey.KeyVaultBaseUrl,
			KeyName:      keyVaultKey.Name,
			KeyVersion:   &keyVaultKey.Version,
			Identity: &factories.CMKIdentityDefinition{
				UserAssignedIdentity: utils.String(d.Get("customer_managed_key_identity_id").(string)),
			},
		}
	}

	globalParameters, err := expandDataFactoryGlobalParameters(d.Get("global_parameter").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	payload.Properties.GlobalParameters = globalParameters

	if _, err := client.CreateOrUpdate(ctx, id, payload, factories.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	githubConfiguration := expandGitHubRepoConfiguration(d.Get("github_configuration").([]interface{}))
	if githubConfiguration != nil {
		repoUpdate := factories.FactoryRepoUpdate{
			FactoryResourceId: utils.String(id.ID()),
			RepoConfiguration: githubConfiguration,
		}
		locationId := factories.NewLocationID(id.SubscriptionId, location)
		if _, err := client.ConfigureFactoryRepo(ctx, locationId, repoUpdate); err != nil {
			return fmt.Errorf("configuring Repository for %s: %+v", locationId, err)
		}
	}
	vstsConfiguration := expandVSTSRepoConfiguration(d.Get("vsts_configuration").([]interface{}))
	if vstsConfiguration != nil {
		repoUpdate := factories.FactoryRepoUpdate{
			FactoryResourceId: utils.String(id.ID()),
			RepoConfiguration: vstsConfiguration,
		}
		locationId := factories.NewLocationID(id.SubscriptionId, location)
		if _, err := client.ConfigureFactoryRepo(ctx, locationId, repoUpdate); err != nil {
			return fmt.Errorf("configuring Repository for %s: %+v", locationId, err)
		}
	}

	if d.Get("managed_virtual_network_enabled").(bool) {
		networkPayload := managedvirtualnetworks.ManagedVirtualNetworkResource{
			Properties: managedvirtualnetworks.ManagedVirtualNetwork{},
		}
		managedNetworkId := managedvirtualnetworks.NewManagedVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName, "default")
		if _, err := managedVirtualNetworksClient.CreateOrUpdate(ctx, managedNetworkId, networkPayload, managedvirtualnetworks.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("creating virtual network for %s: %+v", id, err)
		}
	}

	return resourceDataFactoryRead(d, meta)
}

func resourceDataFactoryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.Factories
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := factories.ParseFactoryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, factories.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FactoryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			customerManagedKeyId := ""
			customerManagedKeyIdentityId := ""
			if enc := props.Encryption; enc != nil {
				if enc.VaultBaseURL != "" && enc.KeyName != "" && enc.KeyVersion != nil {
					version := ""
					if enc.KeyVersion != nil && *enc.KeyVersion != "" {
						version = *enc.KeyVersion
					}
					keyId, err := keyVaultParse.NewNestedKeyID(enc.VaultBaseURL, enc.KeyName, version)
					if err != nil {
						return fmt.Errorf("parsing Nested Item ID: %+v", err)
					}
					customerManagedKeyId = keyId.ID()
				}

				if encIdentity := enc.Identity; encIdentity != nil && encIdentity.UserAssignedIdentity != nil {
					parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*encIdentity.UserAssignedIdentity)
					if err != nil {
						return fmt.Errorf("parsing %q: %+v", *encIdentity.UserAssignedIdentity, err)
					}
					customerManagedKeyIdentityId = parsed.ID()
				}
			}
			d.Set("customer_managed_key_id", customerManagedKeyId)
			d.Set("customer_managed_key_identity_id", customerManagedKeyIdentityId)

			globalParameters, err := flattenDataFactoryGlobalParameters(props.GlobalParameters)
			if err != nil {
				return fmt.Errorf("flattening `global_parameter`: %+v", err)
			}
			if err := d.Set("global_parameter", globalParameters); err != nil {
				return fmt.Errorf("setting `global_parameter`: %+v", err)
			}

			githubConfiguration := flattenGitHubRepoConfiguration(props.RepoConfiguration)
			if err := d.Set("github_configuration", githubConfiguration); err != nil {
				return fmt.Errorf("setting `github_configuration`: %+v", err)
			}

			vstsConfiguration := flattenVSTSRepoConfiguration(props.RepoConfiguration)
			if err := d.Set("vsts_configuration", vstsConfiguration); err != nil {
				return fmt.Errorf("setting `vsts_configuration`: %+v", err)
			}

			// The API defaults this to `true` but won't return it unless it's configured, so default it on
			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil {
				publicNetworkAccessEnabled = *props.PublicNetworkAccess == factories.PublicNetworkAccessEnabled
			}
			d.Set("public_network_enabled", publicNetworkAccessEnabled)

			purviewId := ""
			if purview := props.PurviewConfiguration; purview != nil && purview.PurviewResourceId != nil {
				purviewId = *purview.PurviewResourceId
			}
			d.Set("purview_id", purviewId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	managedVirtualNetworkName, err := getManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, id.SubscriptionId, id.ResourceGroupName, id.FactoryName)
	if err != nil {
		return err
	}
	d.Set("managed_virtual_network_enabled", managedVirtualNetworkName != nil)

	return nil
}

func resourceDataFactoryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.Factories
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := factories.ParseFactoryID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandDataFactoryGlobalParameters(input []interface{}) (*map[string]factories.GlobalParameterSpecification, error) {
	result := make(map[string]factories.GlobalParameterSpecification)
	if len(input) == 0 {
		return &result, nil
	}
	for _, item := range input {
		if item == nil {
			continue
		}
		v := item.(map[string]interface{})

		name := v["name"].(string)
		if _, ok := v[name]; ok {
			return nil, fmt.Errorf("duplicate parameter name")
		}

		result[name] = factories.GlobalParameterSpecification{
			Type:  factories.GlobalParameterType(v["type"].(string)),
			Value: v["value"].(string),
		}
	}
	return &result, nil
}

func flattenDataFactoryGlobalParameters(input *map[string]factories.GlobalParameterSpecification) (*[]interface{}, error) {
	output := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return &output, nil
	}

	for name, item := range *input {
		var valueResult string
		_, valueIsString := item.Value.(string)
		if (item.Type == factories.GlobalParameterTypeArray || item.Type == factories.GlobalParameterTypeObject) && !valueIsString {
			bytes, err := json.Marshal(item.Value)
			if err != nil {
				return nil, fmt.Errorf("marshalling value for global parameter %q (value %+v): %+v", name, item.Value, err)
			}
			valueResult = string(bytes)
		} else {
			valueResult = fmt.Sprintf("%v", item.Value)
		}

		output = append(output, map[string]interface{}{
			"name":  name,
			"type":  string(item.Type),
			"value": valueResult,
		})
	}
	return &output, nil
}

func getManagedVirtualNetworkName(ctx context.Context, client *managedvirtualnetworks.ManagedVirtualNetworksClient, subscriptionId, resourceGroup, factoryName string) (*string, error) {
	factoryId := managedvirtualnetworks.NewFactoryID(subscriptionId, resourceGroup, factoryName)
	resp, err := client.ListByFactory(ctx, factoryId)
	if err != nil {
		return nil, err
	}

	if model := resp.Model; model != nil {
		for _, v := range *model {
			if v.Name == nil {
				continue
			}

			return v.Name, nil
		}
	}

	return nil, nil
}

func expandGitHubRepoConfiguration(input []interface{}) *factories.FactoryGitHubConfiguration {
	if len(input) == 0 {
		return nil
	}

	item := input[0].(map[string]interface{})
	return &factories.FactoryGitHubConfiguration{
		AccountName:         item["account_name"].(string),
		CollaborationBranch: item["branch_name"].(string),
		DisablePublish:      pointer.To(!item["publishing_enabled"].(bool)),
		HostName:            pointer.To(item["git_url"].(string)),
		RepositoryName:      item["repository_name"].(string),
		RootFolder:          item["root_folder"].(string),
	}
}

func flattenGitHubRepoConfiguration(input factories.FactoryRepoConfiguration) []interface{} {
	output := make([]interface{}, 0)

	if v, ok := input.(factories.FactoryGitHubConfiguration); ok {
		gitUrl := ""
		if v.HostName != nil {
			gitUrl = *v.HostName
		}
		publishingEnabled := true
		if v.DisablePublish != nil {
			publishingEnabled = !*v.DisablePublish
		}
		output = append(output, map[string]interface{}{
			"account_name":       v.AccountName,
			"branch_name":        v.CollaborationBranch,
			"git_url":            gitUrl,
			"publishing_enabled": publishingEnabled,
			"repository_name":    v.RepositoryName,
			"root_folder":        v.RootFolder,
		})
	}

	return output
}

func expandVSTSRepoConfiguration(input []interface{}) *factories.FactoryVSTSConfiguration {
	if len(input) == 0 {
		return nil
	}

	item := input[0].(map[string]interface{})
	return &factories.FactoryVSTSConfiguration{
		AccountName:         item["account_name"].(string),
		CollaborationBranch: item["branch_name"].(string),
		DisablePublish:      pointer.To(!item["publishing_enabled"].(bool)),
		ProjectName:         item["project_name"].(string),
		RepositoryName:      item["repository_name"].(string),
		RootFolder:          item["root_folder"].(string),
		TenantId:            pointer.To(item["tenant_id"].(string)),
	}
}

func flattenVSTSRepoConfiguration(input factories.FactoryRepoConfiguration) []interface{} {
	output := make([]interface{}, 0)

	if v, ok := input.(factories.FactoryVSTSConfiguration); ok {
		tenantId := ""
		if v.TenantId != nil {
			tenantId = *v.TenantId
		}
		publishingEnabled := true
		if v.DisablePublish != nil {
			publishingEnabled = !*v.DisablePublish
		}
		output = append(output, map[string]interface{}{
			"account_name":       v.AccountName,
			"branch_name":        v.CollaborationBranch,
			"project_name":       v.ProjectName,
			"publishing_enabled": publishingEnabled,
			"repository_name":    v.RepositoryName,
			"root_folder":        v.RootFolder,
			"tenant_id":          tenantId,
		})
	}

	return output
}
