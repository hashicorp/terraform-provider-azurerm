// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	sentinelmetadata "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MetadataModel struct {
	Name                     string                  `tfschema:"name"`
	WorkspaceId              string                  `tfschema:"workspace_id"`
	ContentId                string                  `tfschema:"content_id"`
	Kind                     string                  `tfschema:"kind"`
	ParentId                 string                  `tfschema:"parent_id"`
	Source                   []MetadataSourceModel   `tfschema:"source"`
	Author                   []MetadataAuthorModel   `tfschema:"author"`
	Support                  []MetadataSupportModel  `tfschema:"support"`
	Dependency               string                  `tfschema:"dependency"`
	Category                 []MetadataCategoryModel `tfschema:"category"`
	Providers                []string                `tfschema:"providers"`
	FirstPublishDate         string                  `tfschema:"first_publish_date"`
	LastPublishDate          string                  `tfschema:"last_publish_date"`
	ContentSchemaVersion     string                  `tfschema:"content_schema_version"`
	CustomVersion            string                  `tfschema:"custom_version"`
	IconId                   string                  `tfschema:"icon_id"`
	PreviewImages            []string                `tfschema:"preview_images"`
	PreviewImagesDark        []string                `tfschema:"preview_images_dark"`
	ThreatAnalysisTactics    []string                `tfschema:"threat_analysis_tactics"`
	ThreatAnalysisTechniques []string                `tfschema:"threat_analysis_techniques"`
	Version                  string                  `tfschema:"version"`
}

type MetadataSourceModel struct {
	Kind string `tfschema:"kind"`
	Name string `tfschema:"name"`
	Id   string `tfschema:"id"`
}

type MetadataAuthorModel struct {
	Name  string `tfschema:"name"`
	Email string `tfschema:"email"`
	Link  string `tfschema:"link"`
}

type MetadataSupportModel struct {
	Tier  string `tfschema:"tier"`
	Name  string `tfschema:"name"`
	Email string `tfschema:"email"`
	Link  string `tfschema:"link"`
}

type MetadataCategoryModel struct {
	Domains   []string `tfschema:"domains"`
	Verticals []string `tfschema:"verticals"`
}
type MetadataResource struct{}

var _ sdk.ResourceWithUpdate = MetadataResource{}

func (a MetadataResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"content_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kind": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(sentinelmetadata.PossibleValuesForKind(), false),
		},

		"parent_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C The API creates a source if omitted but overwriting this/reverting to the default can be done without issue so this can remain
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kind": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(sentinelmetadata.PossibleValuesForSourceKind(), false),
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"author": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"email": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"link": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"support": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"tier": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(sentinelmetadata.SupportTierCommunity),
							string(sentinelmetadata.SupportTierMicrosoft),
							string(sentinelmetadata.SupportTierPartner),
						}, false),
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"email": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"link": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"dependency": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"category": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domains": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"verticals": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"providers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"first_publish_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ISO8601DateTime,
		},

		"last_publish_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ISO8601DateTime,
		},

		"content_schema_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "2.0",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"custom_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"icon_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"preview_images": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"preview_images_dark": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"threat_analysis_tactics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"Reconnaissance",
						"ResourceDevelopment",
						"InitialAccess",
						"Execution",
						"Persistence",
						"PrivilegeEscalation",
						"DefenseEvasion",
						"CredentialAccess",
						"Discovery",
						"LateralMovement",
						"Collection",
						"CommandAndControl",
						"Exfiltration",
						"Impact",
						"ImpairProcessControl",
						"InhibitResponseFunction",
					}, false),
			},
		},

		"threat_analysis_techniques": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (a MetadataResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (a MetadataResource) ModelObject() interface{} {
	return &MetadataModel{}
}

func (a MetadataResource) ResourceType() string {
	return "azurerm_sentinel_metadata"
}

func (a MetadataResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sentinelmetadata.ValidateMetadataID
}

func (a MetadataResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan MetadataModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.Sentinel.MetadataClient

			parsedWorkspaceId, err := workspaces.ParseWorkspaceID(plan.WorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing workspace id: %+v", err)
			}

			id := sentinelmetadata.NewMetadataID(parsedWorkspaceId.SubscriptionId, parsedWorkspaceId.ResourceGroupName, parsedWorkspaceId.WorkspaceName, plan.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(a.ResourceType(), id)
			}

			input := sentinelmetadata.MetadataModel{
				Properties: &sentinelmetadata.MetadataProperties{
					ParentId:   plan.ParentId,
					Kind:       sentinelmetadata.Kind(plan.Kind),
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

			metadata.SetID(id)
			return nil
		},
	}
}

func (a MetadataResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.MetadataClient
			id, err := sentinelmetadata.ParseMetadataID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
			state := MetadataModel{
				Name:        id.MetadataName,
				WorkspaceId: workspaceId.ID(),
			}

			if resp.Model != nil && resp.Model.Properties != nil {
				prop := *resp.Model.Properties

				state.Kind = string(prop.Kind)
				state.ParentId = prop.ParentId
				if prop.ContentId != nil {
					state.ContentId = *prop.ContentId
				}
				if prop.ContentSchemaVersion != nil {
					state.ContentSchemaVersion = *prop.ContentSchemaVersion
				}
				if prop.CustomVersion != nil {
					state.CustomVersion = *prop.CustomVersion
				}
				if prop.FirstPublishDate != nil {
					state.FirstPublishDate = *prop.FirstPublishDate
				}
				if prop.Icon != nil {
					state.IconId = *prop.Icon
				}
				if prop.LastPublishDate != nil {
					state.LastPublishDate = *prop.LastPublishDate
				}
				if prop.PreviewImages != nil {
					state.PreviewImages = *prop.PreviewImages
				}
				if prop.PreviewImagesDark != nil {
					state.PreviewImagesDark = *prop.PreviewImagesDark
				}
				if prop.Providers != nil {
					state.Providers = *prop.Providers
				}
				if prop.ThreatAnalysisTechniques != nil {
					state.ThreatAnalysisTechniques = *prop.ThreatAnalysisTechniques
				}
				if prop.ThreatAnalysisTactics != nil {
					state.ThreatAnalysisTactics = *prop.ThreatAnalysisTactics
				}
				if prop.Version != nil {
					state.Version = *prop.Version
				}

				state.Source = flattenMetadataSourceModel(prop.Source)
				state.Author = flattenMetadataAuthorModel(prop.Author)
				state.Category = flattenMetadataCategoryModel(prop.Categories)
				state.Support = flattenMetadataSupportModel(prop.Support)
				state.Dependency, err = flattenMetadataDependencies(prop.Dependencies)
				if err != nil {
					return fmt.Errorf("flattening `dependency`: %+v", err)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (a MetadataResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.MetadataClient
			id, err := sentinelmetadata.ParseMetadataID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (a MetadataResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan MetadataModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.Sentinel.MetadataClient
			id, err := sentinelmetadata.ParseMetadataID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %q: %+v", metadata.ResourceData.Id(), err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			update := sentinelmetadata.MetadataPatch{
				Properties: &sentinelmetadata.MetadataPropertiesPatch{
					Author:     expandMetadataAuthorModel(plan.Author),
					Categories: expandMetadataCategoryModel(plan.Category),
					Source:     expandMetadataSourceModel(plan.Source),
					Support:    expandMetadataSupportModel(plan.Support),
				},
			}

			if plan.Kind != "" {
				kind := sentinelmetadata.Kind(plan.Kind)
				update.Properties.Kind = &kind
			}

			if plan.ParentId != "" {
				update.Properties.ParentId = &plan.ParentId
			}

			if plan.ContentId != "" {
				update.Properties.ContentId = &plan.ContentId
			}

			if plan.ContentSchemaVersion != "" {
				update.Properties.ContentSchemaVersion = &plan.ContentSchemaVersion
			}

			if plan.CustomVersion != "" {
				update.Properties.CustomVersion = &plan.CustomVersion
			}

			if plan.Dependency != "" {
				depJson, err := pluginsdk.ExpandJsonFromString(plan.Dependency)
				if err != nil {
					return fmt.Errorf("expanding `dependency`: %+v", err)
				}
				update.Properties.Dependencies, err = expandMetadataDependencies(depJson)
				if err != nil {
					return fmt.Errorf("expanding `dependency`: %+v", err)
				}
			}

			if plan.FirstPublishDate != "" {
				update.Properties.FirstPublishDate = &plan.FirstPublishDate
			}

			if plan.IconId != "" {
				update.Properties.Icon = &plan.IconId
			}

			if plan.LastPublishDate != "" {
				update.Properties.LastPublishDate = &plan.LastPublishDate
			}

			if len(plan.PreviewImages) > 0 {
				update.Properties.PreviewImages = &plan.PreviewImages
			}

			if len(plan.PreviewImagesDark) > 0 {
				update.Properties.PreviewImagesDark = &plan.PreviewImagesDark
			}

			if len(plan.Providers) > 0 {
				update.Properties.Providers = &plan.Providers
			}

			if len(plan.ThreatAnalysisTechniques) > 0 {
				update.Properties.ThreatAnalysisTechniques = &plan.ThreatAnalysisTechniques
			}

			if len(plan.ThreatAnalysisTactics) > 0 {
				update.Properties.ThreatAnalysisTactics = &plan.ThreatAnalysisTactics
			}

			if plan.Version != "" {
				update.Properties.Version = &plan.Version
			}

			_, err = client.Update(ctx, *id, update)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandMetadataSourceModel(input []MetadataSourceModel) *sentinelmetadata.MetadataSource {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := sentinelmetadata.MetadataSource{
		Kind: sentinelmetadata.SourceKind(v.Kind),
	}
	if v.Name != "" {
		output.Name = utils.String(v.Name)
	}
	if v.Id != "" {
		output.SourceId = utils.String(v.Id)
	}
	return &output
}

func flattenMetadataSourceModel(input *sentinelmetadata.MetadataSource) []MetadataSourceModel {
	if input == nil {
		return []MetadataSourceModel{}
	}
	output := MetadataSourceModel{
		Kind: string(input.Kind),
	}
	if input.Name != nil {
		output.Name = *input.Name
	}
	if input.SourceId != nil {
		output.Id = *input.SourceId
	}
	return []MetadataSourceModel{output}
}

func expandMetadataAuthorModel(input []MetadataAuthorModel) *sentinelmetadata.MetadataAuthor {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := sentinelmetadata.MetadataAuthor{}
	if v.Name != "" {
		output.Name = utils.String(v.Name)
	}
	if v.Email != "" {
		output.Email = utils.String(v.Email)
	}
	if v.Link != "" {
		output.Link = utils.String(v.Link)
	}
	return &output
}

func flattenMetadataAuthorModel(input *sentinelmetadata.MetadataAuthor) []MetadataAuthorModel {
	if input == nil {
		return []MetadataAuthorModel{}
	}
	output := MetadataAuthorModel{}
	if input.Name != nil {
		output.Name = *input.Name
	}
	if input.Email != nil {
		output.Email = *input.Email
	}
	if input.Link != nil {
		output.Link = *input.Link
	}
	return []MetadataAuthorModel{output}
}

func expandMetadataSupportModel(input []MetadataSupportModel) *sentinelmetadata.MetadataSupport {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := sentinelmetadata.MetadataSupport{
		Tier: sentinelmetadata.SupportTier(v.Tier),
	}
	if v.Name != "" {
		output.Name = utils.String(v.Name)
	}
	if v.Email != "" {
		output.Email = utils.String(v.Email)
	}
	if v.Link != "" {
		output.Link = utils.String(v.Link)
	}
	return &output
}
func flattenMetadataSupportModel(input *sentinelmetadata.MetadataSupport) []MetadataSupportModel {
	if input == nil {
		return []MetadataSupportModel{}
	}
	output := MetadataSupportModel{
		Tier: string(input.Tier),
	}
	if input.Name != nil {
		output.Name = *input.Name
	}
	if input.Email != nil {
		output.Email = *input.Email
	}
	if input.Link != nil {
		output.Link = *input.Link
	}
	return []MetadataSupportModel{output}
}

func expandMetadataCategoryModel(input []MetadataCategoryModel) *sentinelmetadata.MetadataCategories {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	output := sentinelmetadata.MetadataCategories{}
	if len(v.Domains) > 0 {
		output.Domains = &v.Domains
	}
	if len(v.Verticals) > 0 {
		output.Verticals = &v.Verticals
	}
	return &output
}

func flattenMetadataCategoryModel(input *sentinelmetadata.MetadataCategories) []MetadataCategoryModel {
	if input == nil {
		return []MetadataCategoryModel{}
	}
	output := MetadataCategoryModel{}
	if input.Domains != nil {
		output.Domains = *input.Domains
	}
	if input.Verticals != nil {
		output.Verticals = *input.Verticals
	}
	return []MetadataCategoryModel{output}
}

func expandMetadataDependencies(input interface{}) (dependencies *sentinelmetadata.MetadataDependencies, err error) {
	if j, ok := input.(map[string]interface{}); ok {
		dependencies = &sentinelmetadata.MetadataDependencies{}
		// "name" is not returned in response, so it's not supported for now.
		if v, ok := j["contentId"]; ok {
			dependencies.ContentId = utils.String(v.(string))
		}
		if v, ok := j["kind"]; ok {
			kind := sentinelmetadata.Kind(v.(string))
			dependencies.Kind = &kind
		}
		if v, ok := j["version"]; ok {
			dependencies.Version = utils.String(v.(string))
		}
		if v, ok := j["operator"]; ok {
			op := sentinelmetadata.Operator(v.(string))
			dependencies.Operator = &op
		}
		if v, ok := j["criteria"]; ok {
			if array, ok := v.([]interface{}); ok {
				var deps []sentinelmetadata.MetadataDependencies
				for _, item := range array {
					i, ok := item.(map[string]interface{})
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
				dependencies.Criteria = &[]sentinelmetadata.MetadataDependencies{*dep}
			}
		}
		return dependencies, nil
	}

	return nil, fmt.Errorf("unable to parse metadata dependencies: %v", input)
}

func flattenMetadataDependencies(input *sentinelmetadata.MetadataDependencies) (string, error) {
	if input == nil {
		return "", nil
	}
	j, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
