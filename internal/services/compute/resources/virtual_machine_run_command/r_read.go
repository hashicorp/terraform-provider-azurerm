// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_run_command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineRunCommandsClient

			// ErrorBlobManagedIdentity, OutputBlobManagedIdentity, ProtectedParameter, RunAsPassword, Source.ScriptUriManagedIdentity are regarded as sensitive and not returned by API
			var config VirtualMachineRunCommandResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			schema := VirtualMachineRunCommandResourceSchema{
				ErrorBlobManagedIdentity:  config.ErrorBlobManagedIdentity,
				OutputBlobManagedIdentity: config.OutputBlobManagedIdentity,
				ProtectedParameter:        config.ProtectedParameter,
				RunAsPassword:             config.RunAsPassword,
			}

			id, err := virtualmachineruncommands.ParseVirtualMachineRunCommandID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetByVirtualMachine(ctx, *id, virtualmachineruncommands.GetByVirtualMachineOperationOptions{
				// otherwise, the response will not contain instanceView
				Expand: pointer.To("instanceView"),
			})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.Name = id.RunCommandName
			schema.VirtualMachineId = commonids.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineName).ID()

			if model := resp.Model; model != nil {
				schema.Location = model.Location
				schema.Tags = tags.Flatten(model.Tags)
				if prop := model.Properties; prop != nil {
					schema.Parameter = flattenVirtualMachineRunCommandInputParameter(prop.Parameters)
					schema.RunAsUser = pointer.From(prop.RunAsUser)
					schema.InstanceView = flattenVirtualMachineRunCommandInstanceView(prop.InstanceView)
					schema.Source = flattenVirtualMachineRunCommandSource(prop.Source, config)

					// if blob URI is SAS URL, it will not be returned by API
					if strings.Contains(config.ErrorBlobUri, "sig=") {
						schema.ErrorBlobUri = config.ErrorBlobUri
					} else {
						schema.ErrorBlobUri = pointer.From(prop.ErrorBlobUri)
					}

					if strings.Contains(config.OutputBlobUri, "sig=") {
						schema.OutputBlobUri = config.OutputBlobUri
					} else {
						schema.OutputBlobUri = pointer.From(prop.OutputBlobUri)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

type VirtualMachineRunCommandResourceSchema struct {
	ErrorBlobManagedIdentity  []VirtualMachineRunCommandManagedIdentitySchema `tfschema:"error_blob_managed_identity"`
	ErrorBlobUri              string                                          `tfschema:"error_blob_uri"`
	InstanceView              []VirtualMachineRunCommandInstanceViewSchema    `tfschema:"instance_view"`
	Location                  string                                          `tfschema:"location"`
	Name                      string                                          `tfschema:"name"`
	OutputBlobManagedIdentity []VirtualMachineRunCommandManagedIdentitySchema `tfschema:"output_blob_managed_identity"`
	OutputBlobUri             string                                          `tfschema:"output_blob_uri"`
	Parameter                 []VirtualMachineRunCommandInputParameterSchema  `tfschema:"parameter"`
	ProtectedParameter        []VirtualMachineRunCommandInputParameterSchema  `tfschema:"protected_parameter"`
	RunAsPassword             string                                          `tfschema:"run_as_password"`
	RunAsUser                 string                                          `tfschema:"run_as_user"`
	Source                    []VirtualMachineRunCommandScriptSourceSchema    `tfschema:"source"`
	Tags                      map[string]any                                  `tfschema:"tags"`
	VirtualMachineId          string                                          `tfschema:"virtual_machine_id"`
}

type VirtualMachineRunCommandInputParameterSchema struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

type VirtualMachineRunCommandInstanceViewSchema struct {
	ExitCode         int64  `tfschema:"exit_code"`
	executionState   string `tfschema:"execution_state"`
	executionMessage string `tfschema:"execution_message"`
	output           string `tfschema:"output"`
	errorMessage     string `tfschema:"error_message"`
	startTime        string `tfschema:"start_time"`
	endTime          string `tfschema:"end_time"`
}

type VirtualMachineRunCommandManagedIdentitySchema struct {
	ClientId string `tfschema:"client_id"`
	ObjectId string `tfschema:"object_id"`
}

type VirtualMachineRunCommandScriptSourceSchema struct {
	CommandId                string                                          `tfschema:"command_id"`
	Script                   string                                          `tfschema:"script"`
	ScriptUri                string                                          `tfschema:"script_uri"`
	ScriptUriManagedIdentity []VirtualMachineRunCommandManagedIdentitySchema `tfschema:"script_uri_managed_identity"`
}
