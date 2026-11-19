// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_metadata

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	sentinelmetadata "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (a Resource) Read() sdk.ResourceFunc {
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
