// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type LogAnalyticsClusterModel struct {
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Location          string                                     `tfschema:"location"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	SizeGB            int64                                      `tfschema:"size_gb"`
	Tags              map[string]string                          `tfschema:"tags"`
	ClusterId         string                                     `tfschema:"cluster_id"`
}

func (r Resource) ModelObject() any {
	return &LogAnalyticsClusterModel{}
}
