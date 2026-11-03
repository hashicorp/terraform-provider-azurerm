// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_hybrid_connection

type WebAppHybridConnectionModel struct {
	WebAppId            string `tfschema:"web_app_id"`
	RelayId             string `tfschema:"relay_id"`
	HostName            string `tfschema:"hostname"`
	HostPort            int64  `tfschema:"port"`
	SendKeyName         string `tfschema:"send_key_name"`
	NamespaceName       string `tfschema:"namespace_name"`
	RelayName           string `tfschema:"relay_name"`
	ServiceBusNamespace string `tfschema:"service_bus_namespace"`
	ServiceBusSuffix    string `tfschema:"service_bus_suffix"`
	SendKeyValue        string `tfschema:"send_key_value"`
}

func (r Resource) ModelObject() any {
	return &WebAppHybridConnectionModel{}
}
