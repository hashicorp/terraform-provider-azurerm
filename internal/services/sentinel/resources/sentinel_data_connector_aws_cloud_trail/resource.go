// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_data_connector_aws_cloud_trail

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages an AWS CloudTrail Data Connector.",
		Create:      resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Read:        resourceSentinelDataConnectorAwsCloudTrailRead,
		Update:      resourceSentinelDataConnectorAwsCloudTrailCreateUpdate,
		Delete:      resourceSentinelDataConnectorAwsCloudTrailDelete,

		Importer: helpers.ImportDataConnectorUntyped(dataconnectors.DataConnectorKindAmazonWebServicesCloudTrail),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: sentinelDataConnectorAwsCloudTrailSchema(),
	}
}
