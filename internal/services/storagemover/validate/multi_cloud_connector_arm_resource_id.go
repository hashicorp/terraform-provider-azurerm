// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

// MultiCloudConnectorARMResourceID validates a full ARM resource ID for a Hybrid Connectivity public cloud connector.
// See https://learn.microsoft.com/rest/api/storagemover/endpoints/create-or-update
func MultiCloudConnectorARMResourceID(i interface{}, k string) (warnings []string, errors []error) {
	warnings, errors = azure.ValidateResourceID(i, k)
	if len(errors) > 0 {
		return warnings, errors
	}

	id := strings.ToLower(i.(string))
	if !strings.Contains(id, "/providers/microsoft.hybridconnectivity/publiccloudconnectors/") {
		errors = append(errors, fmt.Errorf("%q must be the resource ID of a Microsoft.HybridConnectivity/publicCloudConnectors resource", k))
	}
	return warnings, errors
}

// AwsS3BucketARMResourceID validates a full ARM resource ID for an AWS S3 bucket resource exposed through Microsoft.AwsConnector.
func AwsS3BucketARMResourceID(i interface{}, k string) (warnings []string, errors []error) {
	warnings, errors = azure.ValidateResourceID(i, k)
	if len(errors) > 0 {
		return warnings, errors
	}

	id := strings.ToLower(i.(string))
	if !strings.Contains(id, "/providers/microsoft.awsconnector/s3buckets/") {
		errors = append(errors, fmt.Errorf("%q must be the resource ID of a Microsoft.AwsConnector/s3Buckets resource", k))
	}
	return warnings, errors
}
