// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker_group

type HybridRunbookWorkerGroupModel struct {
	ResourceGroupName     string `tfschema:"resource_group_name"`
	AutomationAccountName string `tfschema:"automation_account_name"`
	Name                  string `tfschema:"name"`
	CredentialName        string `tfschema:"credential_name"`
}

func (m Resource) ModelObject() any {
	return &HybridRunbookWorkerGroupModel{}
}
