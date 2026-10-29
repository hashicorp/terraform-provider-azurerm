// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

type ContainerAppModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroup        string `tfschema:"resource_group_name"`
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`
	Location             string `tfschema:"location"`

	RevisionMode string                      `tfschema:"revision_mode"`
	Ingress      []helpers.Ingress           `tfschema:"ingress"`
	Registries   []helpers.Registry          `tfschema:"registry"`
	Secrets      []helpers.Secret            `tfschema:"secret"`
	Dapr         []helpers.Dapr              `tfschema:"dapr"`
	Template     []helpers.ContainerTemplate `tfschema:"template"`

	Identity             []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	WorkloadProfileName  string                                     `tfschema:"workload_profile_name"`
	MaxInactiveRevisions int64                                      `tfschema:"max_inactive_revisions"`
	Tags                 map[string]any                             `tfschema:"tags"`

	OutboundIpAddresses        []string `tfschema:"outbound_ip_addresses"`
	LatestRevisionName         string   `tfschema:"latest_revision_name"`
	LatestRevisionFqdn         string   `tfschema:"latest_revision_fqdn"`
	CustomDomainVerificationId string   `tfschema:"custom_domain_verification_id"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppModel{}
}
