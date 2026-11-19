// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_recovery_plan

type SiteRecoveryReplicationRecoveryPlanModel struct {
	Name                   string                                         `tfschema:"name"`
	ShutdownRecoveryGroup  []GenericRecoveryGroupModel                    `tfschema:"shutdown_recovery_group"`
	FailoverRecoveryGroup  []GenericRecoveryGroupModel                    `tfschema:"failover_recovery_group"`
	BootRecoveryGroup      []BootRecoveryGroupModel                       `tfschema:"boot_recovery_group"`
	RecoveryVaultId        string                                         `tfschema:"recovery_vault_id"`
	SourceRecoveryFabricId string                                         `tfschema:"source_recovery_fabric_id"`
	TargetRecoveryFabricId string                                         `tfschema:"target_recovery_fabric_id"`
	A2ASettings            []ReplicationRecoveryPlanA2ASpecificInputModel `tfschema:"azure_to_azure_settings"`
}

type GenericRecoveryGroupModel struct {
	PostAction []ActionModel `tfschema:"post_action"`
	PreAction  []ActionModel `tfschema:"pre_action"`
}

type BootRecoveryGroupModel struct {
	PostAction               []ActionModel `tfschema:"post_action"`
	PreAction                []ActionModel `tfschema:"pre_action"`
	ReplicatedProtectedItems []string      `tfschema:"replicated_protected_items"`
}

type ActionModel struct {
	ActionDetailType        string   `tfschema:"type"`
	FabricLocation          string   `tfschema:"fabric_location"`
	FailOverDirections      []string `tfschema:"fail_over_directions"`
	FailOverTypes           []string `tfschema:"fail_over_types"`
	ManualActionInstruction string   `tfschema:"manual_action_instruction"`
	Name                    string   `tfschema:"name"`
	RunbookId               string   `tfschema:"runbook_id"`
	ScriptPath              string   `tfschema:"script_path"`
}

type ReplicationRecoveryPlanA2ASpecificInputModel struct {
	PrimaryZone      string `tfschema:"primary_zone"`
	RecoveryZone     string `tfschema:"recovery_zone"`
	PrimaryEdgeZone  string `tfschema:"primary_edge_zone"`
	RecoveryEdgeZone string `tfschema:"recovery_edge_zone"`
}

func (r Resource) ModelObject() any {
	return &SiteRecoveryReplicationRecoveryPlanModel{}
}
