// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_run_command

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
)

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
			executionState:   pointer.FromEnum(input.ExecutionState),
			executionMessage: pointer.From(input.ExecutionMessage),
			output:           pointer.From(input.Output),
			errorMessage:     pointer.From(input.Error),
			startTime:        pointer.From(input.StartTime),
			endTime:          pointer.From(input.EndTime),
		},
	}
}
