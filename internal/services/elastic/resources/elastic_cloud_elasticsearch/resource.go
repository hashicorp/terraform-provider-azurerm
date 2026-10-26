// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_elasticsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceElasticsearchCreate,
		Read:   resourceElasticsearchRead,
		Update: resourceElasticsearchUpdate,
		Delete: resourceElasticsearchDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := monitorsresource.ParseMonitorID(id)
			return err
		}, resourceElasticsearchImporter),

		Schema: elasticCloudElasticsearchSchema(),
	}
}

func resourceElasticsearchImporter(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
	id, err := elasticmonitorresources.ParseMonitorID(d.Id())
	if err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}

	client := meta.(*clients.Client).Elastic.ServerlessMonitorClient
	resp, err := client.MonitorsGet(ctx, *id)
	if err != nil {
		return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if err := validateElasticCloudHostedMonitor(id, resp.Model); err != nil {
		return []*pluginsdk.ResourceData{d}, err
	}

	return []*pluginsdk.ResourceData{d}, nil
}

func validateElasticCloudHostedMonitor(id *elasticmonitorresources.MonitorId, model *elasticmonitorresources.ElasticMonitorResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}
	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	// Hosted monitors created before hostingType was introduced can omit this field.
	if model.Properties.HostingType == nil {
		return nil
	}
	if *model.Properties.HostingType != elasticmonitorresources.HostingTypeHosted {
		if *model.Properties.HostingType == elasticmonitorresources.HostingTypeServerless {
			return fmt.Errorf("expected %s to use `Hosted` hosting, got `Serverless`; use the `azurerm_elastic_cloud_serverless` resource instead", id)
		}
		return fmt.Errorf("expected %s to use `Hosted` hosting, got `%s`", id, *model.Properties.HostingType)
	}

	return nil
}
