// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana_managed_private_endpoint

type ManagedPrivateEndpointModel struct {
	Name                      string            `tfschema:"name"`
	Location                  string            `tfschema:"location"`
	GrafanaId                 string            `tfschema:"grafana_id"`
	PrivateLinkResourceId     string            `tfschema:"private_link_resource_id"`
	PrivateLinkResourceRegion string            `tfschema:"private_link_resource_region"`
	Tags                      map[string]string `tfschema:"tags"`
	GroupIds                  []string          `tfschema:"group_ids"`
	RequestMessage            string            `tfschema:"request_message"`
	PrivateLinkServiceURL     string            `tfschema:"private_link_service_url"`
}

func (r Resource) ModelObject() any {
	return &ManagedPrivateEndpointModel{}
}
