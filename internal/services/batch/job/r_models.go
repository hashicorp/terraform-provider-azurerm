// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package job

type ResourceModel struct {
	Name                        string            `tfschema:"name"`
	BatchPoolId                 string            `tfschema:"batch_pool_id"`
	DisplayName                 string            `tfschema:"display_name"`
	Priority                    int64             `tfschema:"priority"`
	TaskRetryMaximum            int64             `tfschema:"task_retry_maximum"`
	CommonEnvironmentProperties map[string]string `tfschema:"common_environment_properties"`
}

func (r Resource) ModelObject() interface{} {
	return &ResourceModel{}
}
