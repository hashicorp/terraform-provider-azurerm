// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy_association

// we only support replicate to Azure as backup to customer-managed sites is being deprecated.
// https://learn.microsoft.com/azure/site-recovery/site-to-site-deprecation

type HyperVReplicationPolicyAssociationModel struct {
	Name     string `tfschema:"name"`
	FabricId string `tfschema:"hyperv_site_id"`
	PolicyId string `tfschema:"policy_id"`
}

func (h Resource) ModelObject() any {
	return &HyperVReplicationPolicyAssociationModel{}
}
