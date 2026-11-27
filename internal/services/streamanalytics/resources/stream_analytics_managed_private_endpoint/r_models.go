// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_managed_private_endpoint

type ManagedPrivateEndpointModel struct {
	Name                   string `tfschema:"name"`
	ResourceGroup          string `tfschema:"resource_group_name"`
	StreamAnalyticsCluster string `tfschema:"stream_analytics_cluster_name"`
	TargetResourceId       string `tfschema:"target_resource_id"`
	SubResourceName        string `tfschema:"subresource_name"`
}

func (r Resource) ModelObject() any {
	return &ManagedPrivateEndpointModel{}
}
