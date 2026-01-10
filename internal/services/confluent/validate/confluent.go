// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated - manual changes will be overwritten.

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/connectorresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/organizationresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scclusterrecords"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scenvironmentrecords"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/topicrecords"
)

// OrganizationID validates that the provided ID is a valid Organization ID
var OrganizationID = organizationresources.ValidateOrganizationID

// EnvironmentID validates that the provided ID is a valid Environment ID
var EnvironmentID = scenvironmentrecords.ValidateEnvironmentID

// ClusterID validates that the provided ID is a valid Cluster ID
var ClusterID = scclusterrecords.ValidateClusterID

// TopicID validates that the provided ID is a valid Topic ID
var TopicID = topicrecords.ValidateTopicID

// ConnectorID validates that the provided ID is a valid Connector ID
var ConnectorID = connectorresources.ValidateConnectorID
