// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_hybrid_runbook_worker

type HybridRunbookWorkerModel struct {
	ResourceGroupName     string `tfschema:"resource_group_name"`
	AutomationAccountName string `tfschema:"automation_account_name"`
	WorkerGroupName       string `tfschema:"worker_group_name"`
	WorkerName            string `tfschema:"worker_name"`
	WorkerId              string `tfschema:"worker_id"`
	VmResourceId          string `tfschema:"vm_resource_id"`
	Ip                    string `tfschema:"ip"`
	RegisteredDateTime    string `tfschema:"registration_date_time"`
	LastSeenDateTime      string `tfschema:"last_seen_date_time"`
	WorkerType            string `tfschema:"worker_type"`
}

func (m Resource) ModelObject() any {
	return &HybridRunbookWorkerModel{}
}
