// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_serverless

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
)

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return elasticmonitorresources.ValidateMonitorID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Elastic.ServerlessMonitorClient
		id, err := elasticmonitorresources.ParseMonitorID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := client.MonitorsGet(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return validateElasticCloudServerlessMonitor(id, resp.Model)
	}
}

func (r Resource) Identity() resourceids.ResourceId {
	return &elasticmonitorresources.MonitorId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_elastic_cloud_serverless"
}

func validateElasticCloudServerlessMonitor(id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	if model.Properties.HostingType == nil {
		return fmt.Errorf("retrieving %s: `properties.hostingType` was nil", id)
	}
	if *model.Properties.HostingType != elasticmonitorresources.HostingTypeServerless {
		if *model.Properties.HostingType == elasticmonitorresources.HostingTypeHosted {
			return fmt.Errorf("expected %s to use `Serverless` hosting, got `Hosted`; use the `azurerm_elastic_cloud_elasticsearch` resource instead", id)
		}
		return fmt.Errorf("expected %s to use `Serverless` hosting, got `%s`", id, *model.Properties.HostingType)
	}

	return nil
}
