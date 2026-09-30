// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_run_command

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineRunCommandsClient

			id, err := virtualmachineruncommands.ParseVirtualMachineRunCommandID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetByVirtualMachine(ctx, *id, virtualmachineruncommands.GetByVirtualMachineOperationOptions{
				// otherwise, the response will not contain instanceView
				Expand: pointer.To("instanceView"),
			})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			payload := resp.Model
			if payload.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			var config VirtualMachineRunCommandResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("error_blob_managed_identity") {
				payload.Properties.ErrorBlobManagedIdentity = expandVirtualMachineRunCommandBlobManagedIdentity(config.ErrorBlobManagedIdentity)
			}

			if metadata.ResourceData.HasChange("error_blob_uri") {
				payload.Properties.ErrorBlobUri = pointer.To(config.ErrorBlobUri)
			}

			if metadata.ResourceData.HasChange("output_blob_managed_identity") {
				payload.Properties.OutputBlobManagedIdentity = expandVirtualMachineRunCommandBlobManagedIdentity(config.OutputBlobManagedIdentity)
			}

			if metadata.ResourceData.HasChange("output_blob_uri") {
				payload.Properties.OutputBlobUri = pointer.To(config.OutputBlobUri)
			}

			if metadata.ResourceData.HasChange("parameter") {
				payload.Properties.Parameters = expandVirtualMachineRunCommandInputParameter(config.Parameter)
			}

			if metadata.ResourceData.HasChange("protected_parameter") {
				payload.Properties.ProtectedParameters = expandVirtualMachineRunCommandInputParameter(config.ProtectedParameter)
			}

			if metadata.ResourceData.HasChange("run_as_password") {
				payload.Properties.RunAsPassword = pointer.To(config.RunAsPassword)
			}

			if metadata.ResourceData.HasChange("run_as_user") {
				payload.Properties.RunAsUser = pointer.To(config.RunAsUser)
			}

			if metadata.ResourceData.HasChange("source") {
				payload.Properties.Source = expandVirtualMachineRunCommandSource(config.Source)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
