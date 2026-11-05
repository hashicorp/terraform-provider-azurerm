// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_run_command

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineRunCommandsClient

			var config VirtualMachineRunCommandResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			virtualMachineId, err := commonids.ParseVirtualMachineID(config.VirtualMachineId)
			if err != nil {
				return err
			}

			id := virtualmachineruncommands.NewVirtualMachineRunCommandID(subscriptionId, virtualMachineId.ResourceGroupName, virtualMachineId.VirtualMachineName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetByVirtualMachine(ctx, id, virtualmachineruncommands.DefaultGetByVirtualMachineOperationOptions())
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := virtualmachineruncommands.VirtualMachineRunCommand{
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				Properties: &virtualmachineruncommands.VirtualMachineRunCommandProperties{
					ErrorBlobManagedIdentity:  expandVirtualMachineRunCommandBlobManagedIdentity(config.ErrorBlobManagedIdentity),
					ErrorBlobUri:              pointer.To(config.ErrorBlobUri),
					OutputBlobManagedIdentity: expandVirtualMachineRunCommandBlobManagedIdentity(config.OutputBlobManagedIdentity),
					OutputBlobUri:             pointer.To(config.OutputBlobUri),
					Parameters:                expandVirtualMachineRunCommandInputParameter(config.Parameter),
					ProtectedParameters:       expandVirtualMachineRunCommandInputParameter(config.ProtectedParameter),
					RunAsPassword:             pointer.To(config.RunAsPassword),
					RunAsUser:                 pointer.To(config.RunAsUser),
					Source:                    expandVirtualMachineRunCommandSource(config.Source),

					TimeoutInSeconds: pointer.To(int64(metadata.ResourceData.Timeout(pluginsdk.TimeoutCreate).Seconds())),

					// set API returning error if command run fails
					TreatFailureAsDeploymentFailure: pointer.To(true),
					AsyncExecution:                  pointer.To(false),
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
	}
}

func expandVirtualMachineRunCommandInputParameter(input []VirtualMachineRunCommandInputParameterSchema) *[]virtualmachineruncommands.RunCommandInputParameter {
	output := make([]virtualmachineruncommands.RunCommandInputParameter, 0)

	for _, v := range input {
		parameter := virtualmachineruncommands.RunCommandInputParameter{
			Name:  v.Name,
			Value: v.Value,
		}
		output = append(output, parameter)
	}
	return &output
}

func expandVirtualMachineRunCommandBlobManagedIdentity(input []VirtualMachineRunCommandManagedIdentitySchema) *virtualmachineruncommands.RunCommandManagedIdentity {
	if len(input) == 0 {
		return nil
	}

	output := &virtualmachineruncommands.RunCommandManagedIdentity{}

	if input[0].ClientId != "" {
		output.ClientId = pointer.To(input[0].ClientId)
	}
	if input[0].ObjectId != "" {
		output.ObjectId = pointer.To(input[0].ObjectId)
	}

	return output
}

func expandVirtualMachineRunCommandSource(input []VirtualMachineRunCommandScriptSourceSchema) *virtualmachineruncommands.VirtualMachineRunCommandScriptSource {
	if len(input) == 0 {
		return nil
	}

	output := &virtualmachineruncommands.VirtualMachineRunCommandScriptSource{
		ScriptUriManagedIdentity: expandVirtualMachineRunCommandBlobManagedIdentity(input[0].ScriptUriManagedIdentity),
	}

	if input[0].CommandId != "" {
		output.CommandId = pointer.To(input[0].CommandId)
	}
	if input[0].Script != "" {
		output.Script = pointer.To(input[0].Script)
	}
	if input[0].ScriptUri != "" {
		output.ScriptUri = pointer.To(input[0].ScriptUri)
	}

	return output
}
