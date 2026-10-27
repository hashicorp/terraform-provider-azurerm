// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task_schedule_run_now

type ContainerRegistryTaskScheduleModel struct {
	TaskId string `tfschema:"container_registry_task_id"`
}

func (r Resource) ModelObject() any {
	return &ContainerRegistryTaskScheduleModel{}
}
