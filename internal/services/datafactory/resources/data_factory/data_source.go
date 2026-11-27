// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	datafactoryClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Azure Data Factory (Version 2).",
		Read:        dataSourceDataFactoryRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of this Azure Data Factory.",
				Type:        pluginsdk.TypeString,
				Required:    true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://learn.microsoft.com/azure/data-factory/naming-rules`,
				),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"github_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Description: "The GitHub account name.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"branch_name": {
							Description: "The branch of the repository to get code from.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"git_url": {
							Description: "The GitHub repository URL.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"last_commit_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"repository_name": {
							Description: "The name of the git repository.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"root_folder": {
							Description: "The root folder within the repository.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"vsts_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"account_name": {
							Description: "The VSTS account name.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"branch_name": {
							Description: "The branch of the repository to get code from.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"project_name": {
							Description: "The name of the VSTS project.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"repository_name": {
							Description: "The name of the git repository.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"root_folder": {
							Description: "The root folder within the repository.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"tenant_id": {
							Description: "The Tenant ID associated with the VSTS account.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDataFactoryRead(d *pluginsdk.ResourceData, meta any) error {
	client := datafactoryClient.From(meta).Factories
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := factories.NewFactoryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id, factories.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
			if err := d.Set("github_configuration", flattenGitHubRepoConfigurationDataSource(props.RepoConfiguration)); err != nil {
				return fmt.Errorf("setting `github_configuration`: %+v", err)
			}

			if err := d.Set("vsts_configuration", flattenVSTSRepoConfigurationDataSource(props.RepoConfiguration)); err != nil {
				return fmt.Errorf("setting `vsts_configuration`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func flattenGitHubRepoConfigurationDataSource(input factories.FactoryRepoConfiguration) []any {
	output := make([]any, 0)

	if v, ok := input.(factories.FactoryGitHubConfiguration); ok {
		output = append(output, map[string]any{
			"account_name":    v.AccountName,
			"branch_name":     v.CollaborationBranch,
			"git_url":         pointer.From(v.HostName),
			"last_commit_id":  pointer.From(v.LastCommitId),
			"repository_name": v.RepositoryName,
			"root_folder":     v.RootFolder,
		})
	}

	return output
}

func flattenVSTSRepoConfigurationDataSource(input factories.FactoryRepoConfiguration) []any {
	output := make([]any, 0)

	if v, ok := input.(factories.FactoryVSTSConfiguration); ok {
		output = append(output, map[string]any{
			"account_name":    v.AccountName,
			"branch_name":     v.CollaborationBranch,
			"project_name":    v.ProjectName,
			"repository_name": v.RepositoryName,
			"root_folder":     v.RootFolder,
			"tenant_id":       pointer.From(v.TenantId),
		})
	}

	return output
}
