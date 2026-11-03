// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_service_source_control_slot

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

type SourceControlSlotModel struct {
	SlotID                    string                              `tfschema:"slot_id"`
	SCMType                   string                              `tfschema:"scm_type"`
	RepoURL                   string                              `tfschema:"repo_url"`
	Branch                    string                              `tfschema:"branch"`
	LocalGitSCM               bool                                `tfschema:"use_local_git"`
	ManualIntegration         bool                                `tfschema:"use_manual_integration"`
	UseMercurial              bool                                `tfschema:"use_mercurial"`
	RollbackEnabled           bool                                `tfschema:"rollback_enabled"`
	UsesGithubAction          bool                                `tfschema:"uses_github_action"`
	GithubActionConfiguration []helpers.GithubActionConfiguration `tfschema:"github_action_configuration"`
}

func (r Resource) ModelObject() any {
	return &SourceControlSlotModel{}
}
