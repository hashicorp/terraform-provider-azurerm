// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_subscription_connection

type ManagerSubscriptionConnectionModel struct {
	Name             string `tfschema:"name"`
	SubscriptionId   string `tfschema:"subscription_id"`
	ConnectionState  string `tfschema:"connection_state"`
	Description      string `tfschema:"description"`
	NetworkManagerId string `tfschema:"network_manager_id"`
}

func (r Resource) ModelObject() any {
	return &ManagerSubscriptionConnectionModel{}
}
