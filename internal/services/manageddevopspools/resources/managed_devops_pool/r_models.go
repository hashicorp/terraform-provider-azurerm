// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_devops_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/manageddevopspools/helpers"
)

type ManagedDevOpsPoolModel struct {
	DevCenterProjectId           string                                      `tfschema:"dev_center_project_id"`
	VirtualMachineScaleSetFabric []helpers.VirtualMachineScaleSetFabricModel `tfschema:"virtual_machine_scale_set_fabric"`
	Identity                     []identity.ModelUserAssigned                `tfschema:"identity"`
	Location                     string                                      `tfschema:"location"`
	MaximumConcurrency           int64                                       `tfschema:"maximum_concurrency"`
	Name                         string                                      `tfschema:"name"`
	AzureDevOpsOrganization      []helpers.AzureDevOpsOrganizationModel      `tfschema:"azure_devops_organization"`
	ResourceGroupName            string                                      `tfschema:"resource_group_name"`
	WorkFolder                   string                                      `tfschema:"work_folder"`
	Tags                         map[string]string                           `tfschema:"tags"`
	StatefulAgent                []helpers.StatefulAgentModel                `tfschema:"stateful_agent"`
	StatelessAgent               []helpers.StatelessAgentModel               `tfschema:"stateless_agent"`
}

func (Resource) ModelObject() any {
	return &ManagedDevOpsPoolModel{}
}
