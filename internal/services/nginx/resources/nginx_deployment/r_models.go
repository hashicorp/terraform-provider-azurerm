// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_deployment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type DeploymentModel struct {
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Name                   string                                     `tfschema:"name"`
	NginxVersion           string                                     `tfschema:"nginx_version"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Sku                    string                                     `tfschema:"sku"`
	Location               string                                     `tfschema:"location"`
	Capacity               int64                                      `tfschema:"capacity"`
	AutoScaleProfile       []AutoScaleProfile                         `tfschema:"auto_scale_profile"`
	Email                  string                                     `tfschema:"email"`
	IpAddress              string                                     `tfschema:"ip_address"`
	FrontendPublic         []FrontendPublic                           `tfschema:"frontend_public"`
	FrontendPrivate        []FrontendPrivate                          `tfschema:"frontend_private"`
	NetworkInterface       []NetworkInterface                         `tfschema:"network_interface"`
	UpgradeChannel         string                                     `tfschema:"automatic_upgrade_channel"`
	WebApplicationFirewall []WebApplicationFirewall                   `tfschema:"web_application_firewall"`
	DataplaneAPIEndpoint   string                                     `tfschema:"dataplane_api_endpoint"`
	Tags                   map[string]string                          `tfschema:"tags"`
}

func (m Resource) ModelObject() any {
	return &DeploymentModel{}
}
