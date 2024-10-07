// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = VirtualMachineRunCommandResource{}
	_ sdk.ResourceWithUpdate = VirtualMachineRunCommandResource{}
)

type VirtualMachineRunCommandResource struct{}

func (r VirtualMachineRunCommandResource) ModelObject() interface{} {
	return &VirtualMachineRunCommandResourceSchema{}
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
	Tags                      map[string]interface{}                          `tfschema:"tags"`
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

func (r VirtualMachineRunCommandResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualmachineruncommands.ValidateVirtualMachineRunCommandID
}

func (r VirtualMachineRunCommandResource) ResourceType() string {
	return "azurerm_virtual_machine_run_command"
}

func (r VirtualMachineRunCommandResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VirtualMachineRunCommandName,
		},

		"virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"command_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{
							"source.0.command_id",
							"source.0.script",
							"source.0.script_uri",
						},
					},
					"script": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{
							"source.0.command_id",
							"source.0.script",
							"source.0.script_uri",
						},
					},
					"script_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
						ExactlyOneOf: []string{
							"source.0.command_id",
							"source.0.script",
							"source.0.script_uri",
						},
					},
					"script_uri_managed_identity": {
						Type:      pluginsdk.TypeList,
						Optional:  true,
						Sensitive: true,
						MaxItems:  1,
						RequiredWith: []string{
							"source.0.script_uri",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ConflictsWith: []string{
										"source.0.script_uri_managed_identity.0.object_id",
									},
								},
								"object_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ConflictsWith: []string{
										"source.0.script_uri_managed_identity.0.client_id",
									},
								},
							},
						},
					},
				},
			},
		},

		"error_blob_managed_identity": {
			Type:      pluginsdk.TypeList,
			Optional:  true,
			MaxItems:  1,
			Sensitive: true,
			RequiredWith: []string{
				"error_blob_uri",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsUUID,
						ConflictsWith: []string{
							"error_blob_managed_identity.0.object_id",
						},
					},
					"object_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsUUID,
						ConflictsWith: []string{
							"error_blob_managed_identity.0.client_id",
						},
					},
				},
			},
		},

		"error_blob_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"output_blob_managed_identity": {
			Type:      pluginsdk.TypeList,
			Optional:  true,
			MaxItems:  1,
			Sensitive: true,
			RequiredWith: []string{
				"output_blob_uri",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"output_blob_managed_identity.0.object_id",
						},
					},
					"object_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"output_blob_managed_identity.0.client_id",
						},
					},
				},
			},
		},

		"output_blob_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPS,
		},

		"parameter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"protected_parameter": {
			Type:      pluginsdk.TypeList,
			Optional:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"run_as_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"run_as_user": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r VirtualMachineRunCommandResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"instance_view": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"exit_code": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"execution_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"execution_message": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"output": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"error_message": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"start_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"end_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r VirtualMachineRunCommandResource) Create() sdk.ResourceFunc {
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

			existing, err := client.GetByVirtualMachine(ctx, id, virtualmachineruncommands.DefaultGetByVirtualMachineOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

			result, err := client.CreateOrUpdate(ctx, id, payload)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// the resource still exists if polling fails
			metadata.SetID(id)

			if err := result.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("running the command: %+v", err)
			}

			return nil
		},
	}
}

func (r VirtualMachineRunCommandResource) Read() sdk.ResourceFunc {
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

func (r VirtualMachineRunCommandResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineRunCommandsClient

			id, err := virtualmachineruncommands.ParseVirtualMachineRunCommandID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r VirtualMachineRunCommandResource) Update() sdk.ResourceFunc {
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

func flattenVirtualMachineRunCommandInputParameter(input *[]virtualmachineruncommands.RunCommandInputParameter) []VirtualMachineRunCommandInputParameterSchema {
	if input == nil {
		return make([]VirtualMachineRunCommandInputParameterSchema, 0)
	}

	output := make([]VirtualMachineRunCommandInputParameterSchema, 0)
	for _, v := range *input {
		parameter := VirtualMachineRunCommandInputParameterSchema{
			Name:  v.Name,
			Value: v.Value,
		}
		output = append(output, parameter)

	}

	return output
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

func flattenVirtualMachineRunCommandSource(input *virtualmachineruncommands.VirtualMachineRunCommandScriptSource, config VirtualMachineRunCommandResourceSchema) []VirtualMachineRunCommandScriptSourceSchema {
	if input == nil {
		return []VirtualMachineRunCommandScriptSourceSchema{}
	}

	// if scriptUri is SAS URL, if will not be returned by API
	scriptUri := pointer.From(input.ScriptUri)
	var scriptUriManagedIdentity []VirtualMachineRunCommandManagedIdentitySchema
	if len(config.Source) > 0 {
		if strings.Contains(config.Source[0].ScriptUri, "sig=") {
			scriptUri = config.Source[0].ScriptUri
		}
		scriptUriManagedIdentity = config.Source[0].ScriptUriManagedIdentity
	}

	return []VirtualMachineRunCommandScriptSourceSchema{
		{
			CommandId:                pointer.From(input.CommandId),
			Script:                   pointer.From(input.Script),
			ScriptUri:                scriptUri,
			ScriptUriManagedIdentity: scriptUriManagedIdentity,
		},
	}
}

func flattenVirtualMachineRunCommandInstanceView(input *virtualmachineruncommands.VirtualMachineRunCommandInstanceView) []VirtualMachineRunCommandInstanceViewSchema {
	if input == nil {
		return []VirtualMachineRunCommandInstanceViewSchema{}
	}

	return []VirtualMachineRunCommandInstanceViewSchema{
		{
			ExitCode:         pointer.From(input.ExitCode),
			executionState:   string(pointer.From(input.ExecutionState)),
			executionMessage: pointer.From(input.ExecutionMessage),
			output:           pointer.From(input.Output),
			errorMessage:     pointer.From(input.Error),
			startTime:        pointer.From(input.StartTime),
			endTime:          pointer.From(input.EndTime),
		},
	}
}
