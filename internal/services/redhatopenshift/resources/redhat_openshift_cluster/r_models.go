// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhat_openshift_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type RedHatOpenShiftClusterModel struct {
	Tags                            map[string]string                 `tfschema:"tags"`
	Name                            string                            `tfschema:"name"`
	Location                        string                            `tfschema:"location"`
	ResourceGroup                   string                            `tfschema:"resource_group_name"`
	ConsoleUrl                      string                            `tfschema:"console_url"`
	Identity                        []identity.ModelUserAssigned      `tfschema:"identity"`
	ServicePrincipal                []ServicePrincipal                `tfschema:"service_principal"`
	ClusterProfile                  []ClusterProfile                  `tfschema:"cluster_profile"`
	NetworkProfile                  []NetworkProfile                  `tfschema:"network_profile"`
	MainProfile                     []MainProfile                     `tfschema:"main_profile"`
	WorkerProfile                   []WorkerProfile                   `tfschema:"worker_profile"`
	ApiServerProfile                []ApiServerProfile                `tfschema:"api_server_profile"`
	IngressProfile                  []IngressProfile                  `tfschema:"ingress_profile"`
	PlatformWorkloadIdentityProfile []PlatformWorkloadIdentityProfile `tfschema:"platform_workload_identity_profile"`
}

func (r Resource) ModelObject() any {
	return &RedHatOpenShiftClusterModel{}
}
