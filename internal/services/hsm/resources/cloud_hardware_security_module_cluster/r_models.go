// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cloud_hardware_security_module_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func (Resource) ModelObject() any {
	return &CloudHardwareSecurityModuleClusterModel{}
}

type CloudHardwareSecurityModuleClusterModel struct {
	Name                    string                       `tfschema:"name"`
	ResourceGroupName       string                       `tfschema:"resource_group_name"`
	Location                string                       `tfschema:"location"`
	Identity                []identity.ModelUserAssigned `tfschema:"identity"`
	DomainNameReuse         string                       `tfschema:"domain_name_reuse"`
	Tags                    map[string]string            `tfschema:"tags"`
	HardwareSecurityModules []CloudHsmPropertiesModel    `tfschema:"hardware_security_module"`
}

type CloudHsmPropertiesModel struct {
	FullyQualifiedDomainName string `tfschema:"fully_qualified_domain_name"`
}
