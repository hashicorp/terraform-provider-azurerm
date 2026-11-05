// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_source_control

type SourceControlModel struct {
	AutomationAccountID string     `tfschema:"automation_account_id"`
	Name                string     `tfschema:"name"`
	RepoURL             string     `tfschema:"repository_url"`
	Branch              string     `tfschema:"branch"`
	FolderPath          string     `tfschema:"folder_path"`
	AutoSync            bool       `tfschema:"automatic_sync"`
	PublishRunbook      bool       `tfschema:"publish_runbook_enabled"`
	SourceType          string     `tfschema:"source_control_type"`
	Description         string     `tfschema:"description"`
	SecurityToken       []Security `tfschema:"security"`
}

func (m Resource) ModelObject() any {
	return &SourceControlModel{}
}
