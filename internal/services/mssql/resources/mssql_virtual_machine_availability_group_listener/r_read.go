// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_availability_group_listener

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.VirtualMachinesAvailabilityGroupListenersClient

			id, err := availabilitygrouplisteners.ParseAvailabilityGroupListenerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, availabilitygrouplisteners.GetOperationOptions{Expand: pointer.To("AvailabilityGroupConfiguration")})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := MsSqlVirtualMachineAvailabilityGroupListenerModel{
				Name:                     id.AvailabilityGroupListenerName,
				SqlVirtualMachineGroupId: availabilitygrouplisteners.NewSqlVirtualMachineGroupID(id.SubscriptionId, id.ResourceGroupName, id.SqlVirtualMachineGroupName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.AvailabilityGroupName = pointer.From(props.AvailabilityGroupName)
					state.Port = pointer.From(props.Port)

					avGroupListenerLbConfigs, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(props.LoadBalancerConfigurations, id.SubscriptionId)
					if err != nil {
						return fmt.Errorf("setting `load_balancer_configuration`: %+v", err)
					}
					state.LoadBalancerConfiguration = avGroupListenerLbConfigs

					multiSubnetIpConfiguration, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(props.MultiSubnetIPConfigurations, id.SubscriptionId)
					if err != nil {
						return fmt.Errorf("setting `multi_subnet_ip_configuration`: %+v", err)
					}
					state.MultiSubnetIpConfiguration = multiSubnetIpConfiguration

					if props.AvailabilityGroupConfiguration != nil {
						if props.AvailabilityGroupConfiguration.Replicas != nil {
							replicas, err := flattenMsSqlVirtualMachineAvailabilityGroupListenerReplicas(props.AvailabilityGroupConfiguration.Replicas, id.SubscriptionId)
							if err != nil {
								return fmt.Errorf("setting `replica`: %+v", err)
							}
							state.Replica = replicas
						}
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

type LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener struct {
	LoadBalancerId       string   `tfschema:"load_balancer_id"`
	PrivateIpAddress     string   `tfschema:"private_ip_address"`
	ProbePort            int64    `tfschema:"probe_port"`
	SqlVirtualMachineIds []string `tfschema:"sql_virtual_machine_ids"`
	SubnetId             string   `tfschema:"subnet_id"`
}

type MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener struct {
	PrivateIpAddress    string `tfschema:"private_ip_address"`
	SqlVirtualMachineId string `tfschema:"sql_virtual_machine_id"`
	SubnetId            string `tfschema:"subnet_id"`
}

type ReplicaMsSqlVirtualMachineAvailabilityGroupListener struct {
	SqlVirtualMachineId string `tfschema:"sql_virtual_machine_id"`
	Role                string `tfschema:"role"`
	Commit              string `tfschema:"commit"`
	FailoverMode        string `tfschema:"failover_mode"`
	ReadableSecondary   string `tfschema:"readable_secondary"`
}

func ReplicaSchemaMsSqlVirtualMachineAvailabilityGroupListenerHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", strings.ToLower(m["sql_virtual_machine_id"].(string)))
		fmt.Fprintf(&buf, "%s-", m["role"].(string))
		fmt.Fprintf(&buf, "%s-", m["commit"].(string))
		fmt.Fprintf(&buf, "%s-", m["failover_mode"].(string))
		fmt.Fprintf(&buf, "%s-", m["readable_secondary"].(string))
	}

	return pluginsdk.HashString(buf.String())
}
