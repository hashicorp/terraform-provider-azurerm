// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_manager

func (r Resource) ModelObject() any {
	return &KubernetesFleetManagerResourceSchema{}
}
