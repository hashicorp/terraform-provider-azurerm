// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_socketio

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type WebPubSubSocketIOResourceModel struct {
	Name                             string                                     `tfschema:"name"`
	ResourceGroupName                string                                     `tfschema:"resource_group_name"`
	Location                         string                                     `tfschema:"location"`
	Sku                              []WebPubSubSocketIOSkuModel                `tfschema:"sku"`
	AADAuthEnabled                   bool                                       `tfschema:"aad_auth_enabled"`
	Identity                         []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	LiveTraceEnabled                 bool                                       `tfschema:"live_trace_enabled"`
	LiveTraceConnectivityLogsEnabled bool                                       `tfschema:"live_trace_connectivity_logs_enabled"`
	LiveTraceHttpRequestLogsEnabled  bool                                       `tfschema:"live_trace_http_request_logs_enabled"`
	LiveTraceMessagingLogsEnabled    bool                                       `tfschema:"live_trace_messaging_logs_enabled"`
	LocalAuthEnabled                 bool                                       `tfschema:"local_auth_enabled"`
	PublicNetworkAccess              string                                     `tfschema:"public_network_access"`
	ServiceMode                      string                                     `tfschema:"service_mode"`
	Tags                             map[string]string                          `tfschema:"tags"`
	TlsClientCertEnabled             bool                                       `tfschema:"tls_client_cert_enabled"`
	ExternalIP                       string                                     `tfschema:"external_ip"`
	HostName                         string                                     `tfschema:"hostname"`
	PrimaryAccessKey                 string                                     `tfschema:"primary_access_key"`
	PrimaryConnectionString          string                                     `tfschema:"primary_connection_string"`
	PublicPort                       int64                                      `tfschema:"public_port"`
	ServerPort                       int64                                      `tfschema:"server_port"`
	SecondaryAccessKey               string                                     `tfschema:"secondary_access_key"`
	SecondaryConnectionString        string                                     `tfschema:"secondary_connection_string"`
}

type WebPubSubSocketIOSkuModel struct {
	Name     string `tfschema:"name"`
	Capacity int64  `tfschema:"capacity"`
}

func (w Resource) ModelObject() any {
	return &WebPubSubSocketIOResourceModel{}
}
