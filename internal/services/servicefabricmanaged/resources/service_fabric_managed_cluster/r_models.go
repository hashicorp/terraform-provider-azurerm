// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_managed_cluster

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
)

type ClusterResourceModel struct {
	BackupRestoreService bool   `tfschema:"backup_service_enabled"`
	ClientConnectionPort int64  `tfschema:"client_connection_port"`
	DNSName              string `tfschema:"dns_name"`
	DNSService           bool   `tfschema:"dns_service_enabled"`
	HTTPGatewayPort      int64  `tfschema:"http_gateway_port"`
	Location             string `tfschema:"location"`
	Name                 string `tfschema:"name"`
	Username             string `tfschema:"username"`
	Password             string `tfschema:"password"`
	ResourceGroup        string `tfschema:"resource_group_name"`

	Authentication       []Authentication                     `tfschema:"authentication"`
	CustomFabricSettings []CustomFabricSetting                `tfschema:"custom_fabric_setting"`
	LBRules              []LBRule                             `tfschema:"lb_rule"`
	NodeTypes            []NodeType                           `tfschema:"node_type"`
	Sku                  managedcluster.SkuName               `tfschema:"sku"`
	SubnetId             string                               `tfschema:"subnet_id"`
	Tags                 map[string]any                       `tfschema:"tags"`
	UpgradeWave          managedcluster.ClusterUpgradeCadence `tfschema:"upgrade_wave"`
}

func (k Resource) ModelObject() any {
	return &ClusterResourceModel{}
}
