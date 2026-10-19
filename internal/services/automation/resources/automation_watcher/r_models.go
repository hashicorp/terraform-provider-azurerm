// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_watcher

type WatcherModel struct {
	AutomationAccountID         string         `tfschema:"automation_account_id"`
	Name                        string         `tfschema:"name"`
	Location                    string         `tfschema:"location"`
	Tags                        map[string]any `tfschema:"tags"`
	Etag                        string         `tfschema:"etag"`
	ExecutionFrequencyInSeconds int64          `tfschema:"execution_frequency_in_seconds"`
	ScriptName                  string         `tfschema:"script_name"`
	ScriptParameters            map[string]any `tfschema:"script_parameters"`
	ScriptRunOn                 string         `tfschema:"script_run_on"`
	Description                 string         `tfschema:"description"`
	Status                      string         `tfschema:"status"`
}

func (m Resource) ModelObject() any {
	return &WatcherModel{}
}
