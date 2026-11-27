// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_session_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

type ContainerAppSessionPoolModel struct {
	Name                      string                                     `tfschema:"name"`
	ResourceGroupName         string                                     `tfschema:"resource_group_name"`
	Location                  string                                     `tfschema:"location"`
	ContainerType             string                                     `tfschema:"container_type"`
	MaxConcurrentSessions     int64                                      `tfschema:"max_concurrent_sessions"`
	ContainerAppEnvironmentId string                                     `tfschema:"container_app_environment_id"`
	CustomContainerTemplate   []CustomContainerTemplateModel             `tfschema:"custom_container_template"`
	Identity                  []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	LifecycleConfiguration    []LifecycleConfigurationModel              `tfschema:"lifecycle_configuration"`
	SessionManagedIdentities  []string                                   `tfschema:"session_managed_identities"`
	NetworkEgressEnabled      bool                                       `tfschema:"network_egress_enabled"`
	ReadySessionInstances     int64                                      `tfschema:"ready_session_instances"`
	Secrets                   []SessionPoolSecretModel                   `tfschema:"secret"`
	Tags                      map[string]string                          `tfschema:"tags"`

	NodeCount              int64  `tfschema:"node_count"`
	PoolManagementEndpoint string `tfschema:"pool_management_endpoint"`
}

type LifecycleConfigurationModel struct {
	CooldownPeriodInSeconds int64  `tfschema:"cooldown_period_in_seconds"`
	LifecycleType           string `tfschema:"lifecycle_type"`
	MaxAlivePeriodInSeconds int64  `tfschema:"max_alive_period_in_seconds"`
}

type CustomContainerTemplateModel struct {
	Containers        []SessionContainerModel `tfschema:"container"`
	IngressTargetPort int64                   `tfschema:"ingress_target_port"`
	Registry          []helpers.Registry      `tfschema:"registry"`
}

type SessionContainerModel struct {
	Name    string                    `tfschema:"name"`
	Image   string                    `tfschema:"image"`
	Args    []string                  `tfschema:"args"`
	Command []string                  `tfschema:"command"`
	Cpu     float64                   `tfschema:"cpu"`
	Env     []helpers.ContainerEnvVar `tfschema:"env"`
	Memory  string                    `tfschema:"memory"`
}

type SessionPoolSecretModel struct {
	Name  string `tfschema:"name"`
	Value string `tfschema:"value"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppSessionPoolModel{}
}
