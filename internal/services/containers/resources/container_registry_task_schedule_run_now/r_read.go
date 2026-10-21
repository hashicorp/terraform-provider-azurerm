// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_task_schedule_run_now

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ContainerRegistryTaskScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			model := ContainerRegistryTaskScheduleModel{TaskId: tasks.NewTaskID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TaskName).ID()}
			return metadata.Encode(&model)
		},
	}
}
