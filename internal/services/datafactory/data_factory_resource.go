package datafactory

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			_, err := parse.DataFactoryID(id)
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

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Array",
								"Bool",
								"Float",
								"Int",
								"Object",
								"String",
							}, false),
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

			"tags": tags.Schema(),
		},
	}
}

func resourceDataFactoryCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.FactoriesClient
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDataFactoryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory", id.ID())
		}
	}

	publicNetworkAccess := datafactory.PublicNetworkAccessEnabled
	enabled := d.Get("public_network_enabled").(bool)
	if !enabled {
		publicNetworkAccess = datafactory.PublicNetworkAccessDisabled
	}

	expandedIdentity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	dataFactory := datafactory.Factory{
		Location: utils.String(location),
		FactoryProperties: &datafactory.FactoryProperties{
			PublicNetworkAccess: publicNetworkAccess,
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if purviewId, ok := d.GetOk("purview_id"); ok {
		dataFactory.FactoryProperties.PurviewConfiguration = &datafactory.PurviewConfiguration{
			PurviewResourceID: utils.String(purviewId.(string)),
		}
	}

	if keyVaultKeyID, ok := d.GetOk("customer_managed_key_id"); ok {
		keyVaultKey, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyID.(string))
		if err != nil {
			return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
		}

		dataFactory.FactoryProperties.Encryption = &datafactory.EncryptionConfiguration{
			VaultBaseURL: &keyVaultKey.KeyVaultBaseUrl,
			KeyName:      &keyVaultKey.Name,
			KeyVersion:   &keyVaultKey.Version,
			Identity: &datafactory.CMKIdentityDefinition{
				UserAssignedIdentity: utils.String(d.Get("customer_managed_key_identity_id").(string)),
			},
		}
	}

	globalParameters, err := expandDataFactoryGlobalParameters(d.Get("global_parameter").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	dataFactory.FactoryProperties.GlobalParameters = globalParameters

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, dataFactory, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if hasRepo, repo := expandDataFactoryRepoConfiguration(d); hasRepo {
		repoUpdate := datafactory.FactoryRepoUpdate{
			FactoryResourceID: utils.String(id.ID()),
			RepoConfiguration: repo,
		}
		if _, err := client.ConfigureFactoryRepo(ctx, location, repoUpdate); err != nil {
			return fmt.Errorf("configuring Repository for %s: %+v", id, err)
		}
	}

	managedVirtualNetworkEnabled := d.Get("managed_virtual_network_enabled").(bool)
	// only pass datafactory.ManagedVirtualNetworkResource{} will cause rest api error
	resource := datafactory.ManagedVirtualNetworkResource{
		Properties: &datafactory.ManagedVirtualNetwork{},
	}
	if d.IsNewResource() && managedVirtualNetworkEnabled {
		if _, err := managedVirtualNetworksClient.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, "default", resource, ""); err != nil {
			return fmt.Errorf("creating virtual network for %s: %+v", id, err)
		}
	} else if !d.IsNewResource() && d.HasChange("managed_virtual_network_enabled") {
		if !managedVirtualNetworkEnabled {
			return fmt.Errorf("updating %s: once Managed Virtual Network has been Enabled it's not possible to disable it", id)
		}
		if _, err := managedVirtualNetworksClient.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, "default", resource, ""); err != nil {
			return fmt.Errorf("creating virtual network for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceDataFactoryRead(d, meta)
}

func resourceDataFactoryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.FactoriesClient
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataFactoryID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FactoryName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if factoryProps := resp.FactoryProperties; factoryProps != nil {
		customerManagedKeyId := ""
		customerManagedKeyIdentityId := ""

		if enc := factoryProps.Encryption; enc != nil {
			if enc.VaultBaseURL != nil && enc.KeyName != nil && enc.KeyVersion != nil {
				customerManagedKeyId = fmt.Sprintf("%skeys/%s/%s", *enc.VaultBaseURL, *enc.KeyName, *enc.KeyVersion)
			}

			if enc.Identity != nil && enc.Identity.UserAssignedIdentity != nil {
				parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*enc.Identity.UserAssignedIdentity)
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", *enc.Identity.UserAssignedIdentity, err)
				}
				customerManagedKeyIdentityId = parsed.ID()
			}
		}

		d.Set("customer_managed_key_id", customerManagedKeyId)
		d.Set("customer_managed_key_identity_id", customerManagedKeyIdentityId)

		if err := d.Set("global_parameter", flattenDataFactoryGlobalParameters(factoryProps.GlobalParameters)); err != nil {
			return fmt.Errorf("setting `global_parameter`: %+v", err)
		}
	}

	d.Set("vsts_configuration", []interface{}{})
	d.Set("github_configuration", []interface{}{})
	repoType, repo := flattenDataFactoryRepoConfiguration(&resp)
	if repoType == datafactory.TypeBasicFactoryRepoConfigurationTypeFactoryVSTSConfiguration {
		if err := d.Set("vsts_configuration", repo); err != nil {
			return fmt.Errorf("setting `vsts_configuration`: %+v", err)
		}
	}
	if repoType == datafactory.TypeBasicFactoryRepoConfigurationTypeFactoryGitHubConfiguration {
		if err := d.Set("github_configuration", repo); err != nil {
			return fmt.Errorf("setting `github_configuration`: %+v", err)
		}
	}
	if repoType == datafactory.TypeBasicFactoryRepoConfigurationTypeFactoryRepoConfiguration {
		d.Set("vsts_configuration", repo)
		d.Set("github_configuration", repo)
	}

	identity, err := flattenIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	// This variable isn't returned from the API if it hasn't been passed in first but we know the default is `true`
	if resp.PublicNetworkAccess != "" {
		d.Set("public_network_enabled", resp.PublicNetworkAccess == datafactory.PublicNetworkAccessEnabled)
	}

	if resp.PurviewConfiguration != nil {
		d.Set("purview_id", resp.PurviewConfiguration.PurviewResourceID)
	}

	managedVirtualNetworkEnabled := false
	managedVirtualNetworkName, err := getManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, id.ResourceGroup, id.FactoryName)
	if err != nil {
		return err
	}
	if managedVirtualNetworkName != nil {
		managedVirtualNetworkEnabled = true
	}
	d.Set("managed_virtual_network_enabled", managedVirtualNetworkEnabled)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDataFactoryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.FactoriesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataFactoryID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func expandDataFactoryRepoConfiguration(d *pluginsdk.ResourceData) (bool, datafactory.BasicFactoryRepoConfiguration) {
	if vstsList, ok := d.GetOk("vsts_configuration"); ok {
		vsts := vstsList.([]interface{})[0].(map[string]interface{})
		accountName := vsts["account_name"].(string)
		branchName := vsts["branch_name"].(string)
		projectName := vsts["project_name"].(string)
		repositoryName := vsts["repository_name"].(string)
		rootFolder := vsts["root_folder"].(string)
		tenantID := vsts["tenant_id"].(string)
		return true, &datafactory.FactoryVSTSConfiguration{
			AccountName:         &accountName,
			CollaborationBranch: &branchName,
			ProjectName:         &projectName,
			RepositoryName:      &repositoryName,
			RootFolder:          &rootFolder,
			TenantID:            &tenantID,
		}
	}

	if githubList, ok := d.GetOk("github_configuration"); ok {
		github := githubList.([]interface{})[0].(map[string]interface{})
		accountName := github["account_name"].(string)
		branchName := github["branch_name"].(string)
		gitURL := github["git_url"].(string)
		repositoryName := github["repository_name"].(string)
		rootFolder := github["root_folder"].(string)
		return true, &datafactory.FactoryGitHubConfiguration{
			AccountName:         &accountName,
			CollaborationBranch: &branchName,
			HostName:            &gitURL,
			RepositoryName:      &repositoryName,
			RootFolder:          &rootFolder,
		}
	}

	return false, nil
}

func expandDataFactoryGlobalParameters(input []interface{}) (map[string]*datafactory.GlobalParameterSpecification, error) {
	if len(input) == 0 {
		return nil, nil
	}
	result := make(map[string]*datafactory.GlobalParameterSpecification)
	for _, item := range input {
		if item == nil {
			continue
		}
		v := item.(map[string]interface{})

		name := v["name"].(string)
		if _, ok := v[name]; ok {
			return nil, fmt.Errorf("duplicate parameter name")
		}

		result[name] = &datafactory.GlobalParameterSpecification{
			Type:  datafactory.GlobalParameterType(v["type"].(string)),
			Value: v["value"].(string),
		}
	}
	return result, nil
}

func flattenDataFactoryRepoConfiguration(factory *datafactory.Factory) (datafactory.TypeBasicFactoryRepoConfiguration, []interface{}) {
	result := make([]interface{}, 0)

	if properties := factory.FactoryProperties; properties != nil {
		repo := properties.RepoConfiguration
		if repo != nil {
			settings := map[string]interface{}{}
			if config, test := repo.AsFactoryGitHubConfiguration(); test {
				if config.AccountName != nil {
					settings["account_name"] = *config.AccountName
				}
				if config.CollaborationBranch != nil {
					settings["branch_name"] = *config.CollaborationBranch
				}
				if config.HostName != nil {
					settings["git_url"] = *config.HostName
				}
				if config.RepositoryName != nil {
					settings["repository_name"] = *config.RepositoryName
				}
				if config.RootFolder != nil {
					settings["root_folder"] = *config.RootFolder
				}
				return datafactory.TypeBasicFactoryRepoConfigurationTypeFactoryGitHubConfiguration, append(result, settings)
			}
			if config, test := repo.AsFactoryVSTSConfiguration(); test {
				if config.AccountName != nil {
					settings["account_name"] = *config.AccountName
				}
				if config.CollaborationBranch != nil {
					settings["branch_name"] = *config.CollaborationBranch
				}
				if config.ProjectName != nil {
					settings["project_name"] = *config.ProjectName
				}
				if config.RepositoryName != nil {
					settings["repository_name"] = *config.RepositoryName
				}
				if config.RootFolder != nil {
					settings["root_folder"] = *config.RootFolder
				}
				if config.TenantID != nil {
					settings["tenant_id"] = *config.TenantID
				}
				return datafactory.TypeBasicFactoryRepoConfigurationTypeFactoryVSTSConfiguration, append(result, settings)
			}
		}
	}
	return datafactory.TypeBasicFactoryRepoConfigurationTypeFactoryRepoConfiguration, result
}

func expandIdentity(input []interface{}) (*datafactory.FactoryIdentity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	if expanded.Type == identity.TypeNone {
		return nil, nil
	}

	out := datafactory.FactoryIdentity{
		Type: datafactory.FactoryIdentityType(string(expanded.Type)),
	}

	// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
	if expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.Type = datafactory.FactoryIdentityTypeSystemAssignedUserAssigned
	}
	if len(expanded.IdentityIds) > 0 {
		userAssignedIdentities := make(map[string]interface{})
		for id := range expanded.IdentityIds {
			userAssignedIdentities[id] = make(map[string]interface{})
		}
		out.UserAssignedIdentities = userAssignedIdentities
	}

	return &out, nil
}

func flattenIdentity(input *datafactory.FactoryIdentity) (interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: nil,
		}

		// work around the Swagger defining `SystemAssigned,UserAssigned` rather than `SystemAssigned, UserAssigned`
		if input.Type == datafactory.FactoryIdentityTypeSystemAssignedUserAssigned {
			transform.Type = identity.TypeSystemAssignedUserAssigned
		}

		if input.PrincipalID != nil {
			transform.PrincipalId = input.PrincipalID.String()
		}
		if input.TenantID != nil {
			transform.TenantId = input.TenantID.String()
		}
		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		for k := range input.UserAssignedIdentities {
			identityIds[k] = identity.UserAssignedIdentityDetails{
				// since v is an `interface{}` there's no guarantee this is returned
			}
		}

		transform.IdentityIds = identityIds
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenDataFactoryGlobalParameters(input map[string]*datafactory.GlobalParameterSpecification) []interface{} {
	if len(input) == 0 {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for name, item := range input {
		var valueResult string
		typeResult := azure.TitleCase(string(item.Type))

		if (typeResult == "Array" || typeResult == "Object") && reflect.TypeOf(item.Value).Name() != "string" {
			j, _ := json.Marshal(item.Value)
			valueResult = string(j)
		} else {
			valueResult = fmt.Sprintf("%v", item.Value)
		}

		result = append(result, map[string]interface{}{
			"name":  name,
			"type":  typeResult,
			"value": valueResult,
		})
	}
	return result
}

// Only one VNet is allowed per factory
func getManagedVirtualNetworkName(ctx context.Context, client *datafactory.ManagedVirtualNetworksClient, resourceGroup, factoryName string) (*string, error) {
	resp, err := client.ListByFactory(ctx, resourceGroup, factoryName)
	if err != nil {
		return nil, err
	}
	if len(resp.Values()) == 0 {
		return nil, nil
	}
	managedVirtualNetwork := resp.Values()[0]
	return managedVirtualNetwork.Name, nil
}
