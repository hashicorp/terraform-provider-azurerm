package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactory() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryCreateUpdate,
		Read:   resourceArmDataFactoryRead,
		Update: resourceArmDataFactoryCreateUpdate,
		Delete: resourceArmDataFactoryDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"github_configuration": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"vsts_configuration"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"branch_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"git_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"repository_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"root_folder": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"vsts_configuration": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"github_configuration"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"branch_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"project_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"repository_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"root_folder": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"tenant_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDataFactoryCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.FactoriesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory", *existing.ID)
		}
	}

	dataFactory := datafactory.Factory{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	if v, ok := d.GetOk("identity.0.type"); ok {
		identityType := v.(string)
		dataFactory.Identity = &datafactory.FactoryIdentity{
			Type: &identityType,
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, dataFactory, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory %q (Resource Group %q) ID", name, resourceGroup)
	}

	if hasRepo, repo := expandArmDataFactoryRepoConfiguration(d); hasRepo {
		repoUpdate := datafactory.FactoryRepoUpdate{
			FactoryResourceID: resp.ID,
			RepoConfiguration: repo,
		}
		if _, err = client.ConfigureFactoryRepo(ctx, location, repoUpdate); err != nil {
			return fmt.Errorf("Error configuring Repository for Data Factory %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryRead(d, meta)
}

func resourceArmDataFactoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.FactoriesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["factories"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("vsts_configuration", []interface{}{})
	d.Set("github_configuration", []interface{}{})
	repoType, repo := flattenArmDataFactoryRepoConfiguration(&resp)
	if repoType == datafactory.TypeFactoryVSTSConfiguration {
		if err := d.Set("vsts_configuration", repo); err != nil {
			return fmt.Errorf("Error setting `vsts_configuration`: %+v", err)
		}
	}
	if repoType == datafactory.TypeFactoryGitHubConfiguration {
		if err := d.Set("github_configuration", repo); err != nil {
			return fmt.Errorf("Error setting `github_configuration`: %+v", err)
		}
	}
	if repoType == datafactory.TypeFactoryRepoConfiguration {
		d.Set("vsts_configuration", repo)
		d.Set("github_configuration", repo)
	}

	if err := d.Set("identity", flattenArmDataFactoryIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error flattening `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDataFactoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.FactoriesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["factories"]

	response, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmDataFactoryRepoConfiguration(d *schema.ResourceData) (bool, datafactory.BasicFactoryRepoConfiguration) {
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

func flattenArmDataFactoryRepoConfiguration(factory *datafactory.Factory) (datafactory.TypeBasicFactoryRepoConfiguration, []interface{}) {
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
				return datafactory.TypeFactoryGitHubConfiguration, append(result, settings)
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
				return datafactory.TypeFactoryVSTSConfiguration, append(result, settings)
			}
		}
	}
	return datafactory.TypeFactoryRepoConfiguration, result
}

func flattenArmDataFactoryIdentity(identity *datafactory.FactoryIdentity) interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	if identity.Type != nil {
		result["type"] = *identity.Type
	}
	if identity.PrincipalID != nil {
		result["principal_id"] = identity.PrincipalID.String()
	}
	if identity.TenantID != nil {
		result["tenant_id"] = identity.TenantID.String()
	}

	return []interface{}{result}
}
