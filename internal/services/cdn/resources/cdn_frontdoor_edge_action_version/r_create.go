// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		// deploying the version code can take more than 10 minutes
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.EdgeActionVersionsClient

			var config EdgeActionVersionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			edgeActionId, err := edgeactionversions.ParseEdgeActionID(config.CdnFrontDoorEdgeActionId)
			if err != nil {
				return err
			}

			id := edgeactionversions.NewVersionID(edgeActionId.SubscriptionId, edgeActionId.ResourceGroupName, edgeActionId.EdgeActionName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			// Edge Action Versions are global resources, the API only accepts `global` as the location
			param := edgeactionversions.EdgeActionVersion{
				Location: "global",
				Properties: &edgeactionversions.EdgeActionVersionProperties{
					DeploymentType: edgeactionversions.EdgeActionVersionDeploymentType(config.DeploymentType),
					// the API rejects `True` when another version is already the default, and always makes
					// the first version of an Edge Action the default, so `False` is sent here
					IsDefaultVersion: edgeactionversions.EdgeActionIsDefaultVersionFalse,
				},
				Tags: tags.Expand(config.Tags),
			}

			if err := client.CreateCallbackThenPoll(ctx, id, param, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			if config.CodeContentBase64 != "" {
				if err := deployVersionCode(ctx, client, id, config); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
