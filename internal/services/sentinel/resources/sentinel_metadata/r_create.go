// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_metadata

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (a Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			var plan MetadataModel
			if err := meta.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := meta.Client.Sentinel.MetadataClient

			parsedWorkspaceId, err := workspaces.ParseWorkspaceID(plan.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			id := metadata.NewMetadataID(parsedWorkspaceId.SubscriptionId, parsedWorkspaceId.ResourceGroupName, parsedWorkspaceId.WorkspaceName, plan.Name)

			if !meta.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return meta.ResourceRequiresImport(a.ResourceType(), id)
				}
			}

			input := metadata.MetadataModel{
				Properties: &metadata.MetadataProperties{
					ParentId:   plan.ParentId,
					Kind:       metadata.Kind(plan.Kind),
					Author:     expandMetadataAuthorModel(plan.Author),
					Categories: expandMetadataCategoryModel(plan.Category),
					Source:     expandMetadataSourceModel(plan.Source),
					Support:    expandMetadataSupportModel(plan.Support),
				},
			}

			if plan.ContentId != "" {
				input.Properties.ContentId = &plan.ContentId
			}

			if plan.ContentSchemaVersion != "" {
				input.Properties.ContentSchemaVersion = &plan.ContentSchemaVersion
			}

			if plan.CustomVersion != "" {
				input.Properties.CustomVersion = &plan.CustomVersion
			}

			if plan.Dependency != "" {
				depJson, err := pluginsdk.ExpandJsonFromString(plan.Dependency)
				if err != nil {
					return fmt.Errorf("expanding `dependency`: %+v", err)
				}
				input.Properties.Dependencies, err = expandMetadataDependencies(depJson)
				if err != nil {
					return fmt.Errorf("expanding `dependency`: %+v", err)
				}
			}

			if plan.FirstPublishDate != "" {
				input.Properties.FirstPublishDate = &plan.FirstPublishDate
			}

			if plan.IconId != "" {
				input.Properties.Icon = &plan.IconId
			}

			if plan.LastPublishDate != "" {
				input.Properties.LastPublishDate = &plan.LastPublishDate
			}

			if len(plan.PreviewImages) > 0 {
				input.Properties.PreviewImages = &plan.PreviewImages
			}

			if len(plan.PreviewImagesDark) > 0 {
				input.Properties.PreviewImagesDark = &plan.PreviewImagesDark
			}

			if len(plan.Providers) > 0 {
				input.Properties.Providers = &plan.Providers
			}

			if len(plan.ThreatAnalysisTechniques) > 0 {
				input.Properties.ThreatAnalysisTechniques = &plan.ThreatAnalysisTechniques
			}

			if len(plan.ThreatAnalysisTactics) > 0 {
				input.Properties.ThreatAnalysisTactics = &plan.ThreatAnalysisTactics
			}

			if plan.Version != "" {
				input.Properties.Version = &plan.Version
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func expandMetadataSourceModel(input []MetadataSourceModel) *metadata.MetadataSource {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := metadata.MetadataSource{
		Kind: metadata.SourceKind(v.Kind),
	}
	if v.Name != "" {
		output.Name = pointer.To(v.Name)
	}
	if v.Id != "" {
		output.SourceId = pointer.To(v.Id)
	}
	return &output
}

func expandMetadataAuthorModel(input []MetadataAuthorModel) *metadata.MetadataAuthor {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := metadata.MetadataAuthor{}
	if v.Name != "" {
		output.Name = pointer.To(v.Name)
	}
	if v.Email != "" {
		output.Email = pointer.To(v.Email)
	}
	if v.Link != "" {
		output.Link = pointer.To(v.Link)
	}
	return &output
}

func expandMetadataSupportModel(input []MetadataSupportModel) *metadata.MetadataSupport {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := metadata.MetadataSupport{
		Tier: metadata.SupportTier(v.Tier),
	}
	if v.Name != "" {
		output.Name = pointer.To(v.Name)
	}
	if v.Email != "" {
		output.Email = pointer.To(v.Email)
	}
	if v.Link != "" {
		output.Link = pointer.To(v.Link)
	}
	return &output
}

func expandMetadataCategoryModel(input []MetadataCategoryModel) *metadata.MetadataCategories {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := metadata.MetadataCategories{}
	if len(v.Domains) > 0 {
		output.Domains = &v.Domains
	}
	if len(v.Verticals) > 0 {
		output.Verticals = &v.Verticals
	}
	return &output
}

func expandMetadataDependencies(input any) (dependencies *metadata.MetadataDependencies, err error) {
	if j, ok := input.(map[string]any); ok {
		dependencies = &metadata.MetadataDependencies{}
		// "name" is not returned in response, so it's not supported for now.
		if v, ok := j["contentId"]; ok {
			dependencies.ContentId = pointer.To(v.(string))
		}
		if v, ok := j["kind"]; ok {
			dependencies.Kind = pointer.ToEnum[metadata.Kind](v.(string))
		}
		if v, ok := j["version"]; ok {
			dependencies.Version = pointer.To(v.(string))
		}
		if v, ok := j["operator"]; ok {
			dependencies.Operator = pointer.ToEnum[metadata.Operator](v.(string))
		}
		if v, ok := j["criteria"]; ok {
			if array, ok := v.([]any); ok {
				var deps []metadata.MetadataDependencies
				for _, item := range array {
					i, ok := item.(map[string]any)
					if !ok {
						continue
					}
					if len(i) == 0 {
						continue
					}
					dep, err := expandMetadataDependencies(i)
					if err != nil {
						return nil, err
					}
					deps = append(deps, *dep)
				}
				dependencies.Criteria = &deps
			} else {
				dep, err := expandMetadataDependencies(v)
				if err != nil {
					return nil, err
				}
				dependencies.Criteria = &[]metadata.MetadataDependencies{*dep}
			}
		}
		return dependencies, nil
	}

	return nil, fmt.Errorf("unable to parse metadata dependencies: %v", input)
}
