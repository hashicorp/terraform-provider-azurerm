// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_auto_upgrade_profile

func (r Resource) ModelObject() any {
	return &KubernetesFleetAutoUpgradeProfileResourceModel{}
}

type KubernetesFleetAutoUpgradeProfileResourceModel struct {
	Name                     string `tfschema:"name"`
	KubernetesFleetManagerId string `tfschema:"kubernetes_fleet_manager_id"`
	Channel                  string `tfschema:"channel"`
	NodeImageSelectionType   string `tfschema:"node_image_selection_type"`
	UpdateStrategyId         string `tfschema:"update_strategy_id"`
	Enabled                  bool   `tfschema:"enabled"`
}
