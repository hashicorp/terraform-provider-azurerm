// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_container_immutability_policy

type ContainerImmutabilityPolicyModel struct {
	StorageContainerResourceManagerId string `tfschema:"storage_container_resource_manager_id"`
	ImmutabilityPeriodInDays          int64  `tfschema:"immutability_period_in_days"`
	Locked                            bool   `tfschema:"locked"`
	ProtectedAppendWritesAllEnabled   bool   `tfschema:"protected_append_writes_all_enabled"`
	ProtectedAppendWritesEnabled      bool   `tfschema:"protected_append_writes_enabled"`
}

func (r Resource) ModelObject() any {
	return &ContainerImmutabilityPolicyModel{}
}
