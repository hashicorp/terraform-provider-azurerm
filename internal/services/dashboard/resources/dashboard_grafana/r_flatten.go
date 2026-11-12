// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenSMTPConfigurationModel(input *managedgrafanas.Smtp, data *schema.ResourceData) []SMTPConfigurationModel {
	var outputList []SMTPConfigurationModel
	if input == nil || !pointer.From(input.Enabled) {
		return []SMTPConfigurationModel{}
	}

	var output SMTPConfigurationModel

	if input.Enabled != nil {
		output.SMTPEnabled = pointer.From(input.Enabled)
	}

	if input.FromAddress != nil {
		output.FromAddress = pointer.From(input.FromAddress)
	}

	if input.FromName != nil {
		output.FromName = pointer.From(input.FromName)
	}

	if input.Host != nil {
		output.Host = pointer.From(input.Host)
	}

	if input.SkipVerify != nil {
		output.SkipVerify = pointer.From(input.SkipVerify)
	}

	if input.User != nil {
		output.User = pointer.From(input.User)
	}

	if input.StartTLSPolicy != nil {
		output.StartTLSPolicy = pointer.FromEnum(input.StartTLSPolicy)
	}

	output.Password = data.Get("smtp.0.password").(string)

	return append(outputList, output)
}

func flattenAzureMonitorWorkspaceIntegrationModelArray(inputList *[]managedgrafanas.AzureMonitorWorkspaceIntegration) []AzureMonitorWorkspaceIntegrationModel {
	if inputList == nil {
		return []AzureMonitorWorkspaceIntegrationModel{}
	}

	outputList := make([]AzureMonitorWorkspaceIntegrationModel, 0, len(*inputList))
	for _, input := range *inputList {
		output := AzureMonitorWorkspaceIntegrationModel{}

		if input.AzureMonitorWorkspaceResourceId != nil {
			output.ResourceId = *input.AzureMonitorWorkspaceResourceId
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenLegacySystemAndUserAssignedMap(input *identity.LegacySystemAndUserAssignedMap) *[]any {
	if input == nil {
		return &[]any{}
	}

	identityValue := &identity.SystemOrUserAssignedMap{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
		IdentityIds: input.IdentityIds,
	}

	output, err := identity.FlattenSystemOrUserAssignedMap(identityValue)
	if err != nil {
		return &[]any{}
	}
	return output
}
