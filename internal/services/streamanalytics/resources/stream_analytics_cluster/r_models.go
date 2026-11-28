// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_cluster

type ClusterModel struct {
	Name              string         `tfschema:"name"`
	ResourceGroup     string         `tfschema:"resource_group_name"`
	Location          string         `tfschema:"location"`
	StreamingCapacity int64          `tfschema:"streaming_capacity"`
	Tags              map[string]any `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ClusterModel{}
}
