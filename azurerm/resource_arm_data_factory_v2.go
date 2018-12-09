package azurerm

import (
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryV2CreateOrUpdate,
		Read:   resourceArmDataFactoryV2Read,
		Update: resourceArmDataFactoryV2CreateOrUpdate,
		Delete: resourceArmDataFactoryV2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules`,
				),
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, true),
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"collaboration_branch": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repository_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"root_folder": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{"vsts_configuration"},
			},

			"vsts_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"collaboration_branch": {
							Type:     schema.TypeString,
							Required: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repository_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"root_folder": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{"github_configuration"},
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDataFactoryV2CreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	dataFactory := datafactory.Factory{
		Location: &location,
		Tags:     expandTags(tags),
	}

	if v, ok := d.GetOk("identity.0.type"); ok {
		identityType := v.(string)
		dataFactory.Identity = &datafactory.FactoryIdentity{
			Type: &identityType,
		}
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, name, dataFactory, "")
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory %s (resource group %s) ID", name, resourceGroup)
	}

	if hasRepo, repo := expandArmDataFactoryV2RepoConfiguration(d); hasRepo {
		repoUpdate := datafactory.FactoryRepoUpdate{
			FactoryResourceID: resp.ID,
			RepoConfiguration: repo,
		}
		_, err = client.ConfigureFactoryRepo(ctx, location, repoUpdate)
		if err != nil {
			return err
		}
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryV2Read(d, meta)
}

func resourceArmDataFactoryV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

		return fmt.Errorf("Error making Read request on AzureRM Data Factory %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	repoType, repo := flattenArmDataFactoryV2RepoConfiguration(&resp)
	if repoType == datafactory.TypeFactoryVSTSConfiguration {
		if err := d.Set("vsts_configuration", repo); err != nil {
			return fmt.Errorf("Error setting Data Factory %q `vsts_configuration` (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if repoType == datafactory.TypeFactoryGitHubConfiguration {
		if err := d.Set("github_configuration", repo); err != nil {
			return fmt.Errorf("Error setting Data Factory %q `github_configuration` (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if repoType == datafactory.TypeFactoryRepoConfiguration {
		d.Set("vsts_configuration", repo)
		d.Set("github_configuration", repo)
	}

	if err := d.Set("identity", flattenArmDataFactoryV2Identity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting Data Factory %q `identity` (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if tags := resp.Tags; tags != nil {
		flattenAndSetTags(d, tags)
	}

	return nil
}

func resourceArmDataFactoryV2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["factories"]

	response, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory %s: %+v", name, err)
		}
	}

	return nil
}

func expandArmDataFactoryV2RepoConfiguration(d *schema.ResourceData) (bool, datafactory.BasicFactoryRepoConfiguration) {
	if vstsList, ok := d.GetOk("vsts_configuration"); ok {
		vsts := vstsList.([]interface{})[0].(map[string]interface{})
		accountName := vsts["account_name"].(string)
		collaborationBranch := vsts["collaboration_branch"].(string)
		projectName := vsts["project_name"].(string)
		repositoryName := vsts["repository_name"].(string)
		rootFolder := vsts["root_folder"].(string)
		tenantID := vsts["tenant_id"].(string)
		// https://github.com/Azure/go-autorest/issues/307
		return true, &datafactory.FactoryVSTSConfiguration{
			AccountName:         &accountName,
			CollaborationBranch: &collaborationBranch,
			ProjectName:         &projectName,
			RepositoryName:      &repositoryName,
			RootFolder:          &rootFolder,
			TenantID:            &tenantID,
		}
	}

	if githubList, ok := d.GetOk("github_configuration"); ok {
		github := githubList.([]interface{})[0].(map[string]interface{})
		accountName := github["account_name"].(string)
		collaborationBranch := github["collaboration_branch"].(string)
		hostName := github["host_name"].(string)
		repositoryName := github["repository_name"].(string)
		rootFolder := github["root_folder"].(string)
		// https://github.com/Azure/go-autorest/issues/307
		return true, &datafactory.FactoryGitHubConfiguration{
			AccountName:         &accountName,
			CollaborationBranch: &collaborationBranch,
			HostName:            &hostName,
			RepositoryName:      &repositoryName,
			RootFolder:          &rootFolder,
		}
	}

	return false, nil
}

func flattenArmDataFactoryV2RepoConfiguration(factory *datafactory.Factory) (datafactory.TypeBasicFactoryRepoConfiguration, []interface{}) {
	result := make([]interface{}, 0)
	properties := factory.FactoryProperties
	if properties != nil {
		repo := properties.RepoConfiguration
		if repo != nil {
			settings := map[string]interface{}{}
			if config, test := repo.AsFactoryGitHubConfiguration(); test {
				settings["account_name"] = *config.AccountName
				settings["collaboration_branch"] = *config.CollaborationBranch
				settings["host_name"] = *config.HostName
				settings["repository_name"] = *config.RepositoryName
				settings["root_folder"] = *config.RootFolder
				return datafactory.TypeFactoryGitHubConfiguration, append(result, settings)
			}
			if config, test := repo.AsFactoryVSTSConfiguration(); test {
				settings["account_name"] = *config.AccountName
				settings["collaboration_branch"] = *config.CollaborationBranch
				settings["project_name"] = *config.ProjectName
				settings["repository_name"] = *config.RepositoryName
				settings["root_folder"] = *config.RootFolder
				settings["tenant_id"] = *config.TenantID
				return datafactory.TypeFactoryVSTSConfiguration, append(result, settings)
			}
		}
	}
	return datafactory.TypeFactoryRepoConfiguration, result
}

func flattenArmDataFactoryV2Identity(identity *datafactory.FactoryIdentity) interface{} {
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
