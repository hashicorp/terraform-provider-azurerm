// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_metadata

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	sentinelmetadata "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (a Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan MetadataModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := sentinelClient.From(metadata.Client).MetadataClient
			id, err := sentinelmetadata.ParseMetadataID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %q: %+v", metadata.ResourceData.Id(), err)
			}

			if _, err = client.Get(ctx, *id); err != nil {
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
				update.Properties.Kind = pointer.ToEnum[sentinelmetadata.Kind](plan.Kind)
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

			if _, err = client.Update(ctx, *id, update); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
