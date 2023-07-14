// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDataFactory() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataFactoryRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://docs.microsoft.com/azure/data-factory/naming-rules`,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"branch_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"git_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"repository_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"root_folder": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"branch_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"repository_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"root_folder": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDataFactoryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.Factories
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
			githubConfiguration := flattenGitHubRepoConfigurationDataSource(props.RepoConfiguration)
			if err := d.Set("github_configuration", githubConfiguration); err != nil {
				return fmt.Errorf("setting `github_configuration`: %+v", err)
			}

			vstsConfiguration := flattenVSTSRepoConfigurationDataSource(props.RepoConfiguration)
			if err := d.Set("vsts_configuration", vstsConfiguration); err != nil {
				return fmt.Errorf("setting `vsts_configuration`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func flattenGitHubRepoConfigurationDataSource(input factories.FactoryRepoConfiguration) []interface{} {
	output := make([]interface{}, 0)

	if v, ok := input.(factories.FactoryGitHubConfiguration); ok {
		gitUrl := ""
		if v.HostName != nil {
			gitUrl = *v.HostName
		}
		output = append(output, map[string]interface{}{
			"account_name":    v.AccountName,
			"branch_name":     v.CollaborationBranch,
			"git_url":         gitUrl,
			"repository_name": v.RepositoryName,
			"root_folder":     v.RootFolder,
		})
	}

	return output
}

func flattenVSTSRepoConfigurationDataSource(input factories.FactoryRepoConfiguration) []interface{} {
	output := make([]interface{}, 0)

	if v, ok := input.(factories.FactoryVSTSConfiguration); ok {
		tenantId := ""
		if v.TenantId != nil {
			tenantId = *v.TenantId
		}
		output = append(output, map[string]interface{}{
			"account_name":    v.AccountName,
			"branch_name":     v.CollaborationBranch,
			"project_name":    v.ProjectName,
			"repository_name": v.RepositoryName,
			"root_folder":     v.RootFolder,
			"tenant_id":       tenantId,
		})
	}

	return output
}
