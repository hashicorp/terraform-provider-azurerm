// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_elasticsearch

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
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

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := monitorsresource.ParseMonitorID(id)
			return err
		}),

		Schema: elasticCloudElasticsearchSchema(),
	}
}
