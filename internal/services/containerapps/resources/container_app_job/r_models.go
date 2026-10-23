// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_job

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

type ContainerAppJobModel struct {
	Name                      string                                     `tfschema:"name"`
	ResourceGroup             string                                     `tfschema:"resource_group_name"`
	Location                  string                                     `tfschema:"location"`
	ContainerAppEnvironmentId string                                     `tfschema:"container_app_environment_id"`
	WorkloadProfileName       string                                     `tfschema:"workload_profile_name"`
	Template                  []helpers.JobTemplateModel                 `tfschema:"template"`
	ReplicaRetryLimit         int64                                      `tfschema:"replica_retry_limit"`
	ReplicaTimeoutInSeconds   int64                                      `tfschema:"replica_timeout_in_seconds"`
	Secrets                   []helpers.Secret                           `tfschema:"secret"`
	Registries                []helpers.Registry                         `tfschema:"registry"`
	EventTriggerConfig        []helpers.EventTriggerConfiguration        `tfschema:"event_trigger_config"`
	ManualTriggerConfig       []helpers.ManualTriggerConfiguration       `tfschema:"manual_trigger_config"`
	ScheduleTriggerConfig     []helpers.ScheduleTriggerConfiguration     `tfschema:"schedule_trigger_config"`
	Identity                  []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags                      map[string]any                             `tfschema:"tags"`

	OutboundIPAddresses []string `tfschema:"outbound_ip_addresses"`
	EventStreamEndpoint string   `tfschema:"event_stream_endpoint"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppJobModel{}
}
