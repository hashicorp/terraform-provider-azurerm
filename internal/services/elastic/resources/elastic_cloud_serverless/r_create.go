// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_serverless

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ElasticCloudServerlessResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Elastic.ServerlessMonitorClient
			id := elasticmonitorresources.NewMonitorID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MonitorsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			monitoringStatus := elasticmonitorresources.MonitoringStatusDisabled
			if model.MonitoringEnabled {
				monitoringStatus = elasticmonitorresources.MonitoringStatusEnabled
			}

			body := elasticmonitorresources.ElasticMonitorResource{
				Kind:     pointer.To(model.Kind),
				Location: location.Normalize(model.Location),
				Properties: &elasticmonitorresources.MonitorProperties{
					GenerateApiKey:   pointer.To(model.GenerateAPIKey),
					HostingType:      pointer.To(elasticmonitorresources.HostingTypeServerless),
					MonitoringStatus: pointer.To(monitoringStatus),
					PlanDetails: &elasticmonitorresources.PlanDetails{
						OfferID:     pointer.To(model.OfferID),
						PlanID:      pointer.To(model.PlanID),
						PublisherID: pointer.To(model.PublisherID),
						TermID:      pointer.To(model.TermID),
					},
					ProjectDetails: &elasticmonitorresources.ProjectDetails{
						ConfigurationType: pointer.ToEnum[elasticmonitorresources.ConfigurationType](model.ConfigurationType),
						ProjectType:       pointer.ToEnum[elasticmonitorresources.ProjectType](model.ProjectType),
					},
					UserInfo: &elasticmonitorresources.UserInfo{
						EmailAddress: pointer.To(model.ElasticCloudEmailAddress),
					},
				},
				Sku: &elasticmonitorresources.ResourceSku{
					Name: model.Sku,
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.MonitorsCreateCallbackThenPoll(ctx, id, body, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}
