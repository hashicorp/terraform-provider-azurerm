// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_availability_group_listener

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlVirtualMachineAvailabilityGroupListenerModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.MSSQL.VirtualMachinesAvailabilityGroupListenersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			sqlVirtualMachineGroupId, err := availabilitygrouplisteners.ParseSqlVirtualMachineGroupID(model.SqlVirtualMachineGroupId)
			if err != nil {
				return err
			}

			id := availabilitygrouplisteners.NewAvailabilityGroupListenerID(subscriptionId, sqlVirtualMachineGroupId.ResourceGroupName, sqlVirtualMachineGroupId.SqlVirtualMachineGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, availabilitygrouplisteners.GetOperationOptions{Expand: pointer.To("AvailabilityGroupConfiguration")})
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			replicas, err := expandMsSqlVirtualMachineAvailabilityGroupListenerReplicas(model.Replica)
			if err != nil {
				return err
			}

			parameters := availabilitygrouplisteners.AvailabilityGroupListener{
				Properties: &availabilitygrouplisteners.AvailabilityGroupListenerProperties{
					AvailabilityGroupName:                    pointer.To(model.AvailabilityGroupName),
					CreateDefaultAvailabilityGroupIfNotExist: pointer.To(true),
					Port:                                     pointer.To(model.Port),
					AvailabilityGroupConfiguration: &availabilitygrouplisteners.AgConfiguration{
						Replicas: replicas,
					},
				},
			}

			if len(model.LoadBalancerConfiguration) != 0 {
				lbConfigs, err := expandMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(model.LoadBalancerConfiguration)
				if err != nil {
					return err
				}
				parameters.Properties.LoadBalancerConfigurations = lbConfigs
			}

			if len(model.MultiSubnetIpConfiguration) != 0 {
				parameters.Properties.MultiSubnetIPConfigurations = expandMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(model.MultiSubnetIpConfiguration)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(lbConfigs []LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener) (*[]availabilitygrouplisteners.LoadBalancerConfiguration, error) {
	results := make([]availabilitygrouplisteners.LoadBalancerConfiguration, 0)

	for _, lb := range lbConfigs {
		lbConfig := availabilitygrouplisteners.LoadBalancerConfiguration{
			ProbePort: pointer.To(lb.ProbePort),
		}

		parsedLbId := ""
		if lb.LoadBalancerId != "" {
			id, err := loadbalancers.ParseLoadBalancerID(lb.LoadBalancerId)
			if err != nil {
				return nil, err
			}
			parsedLbId = id.ID()
		}
		lbConfig.LoadBalancerResourceId = pointer.To(parsedLbId)

		var parsedIds []any
		for _, sqlVmId := range lb.SqlVirtualMachineIds {
			parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineID(sqlVmId)
			if err != nil {
				return nil, err
			}
			parsedIds = append(parsedIds, parsedId.ID())
		}
		lbConfig.SqlVirtualMachineInstances = pluginsdk.ExpandStringSlice(parsedIds)

		lbConfig.PrivateIPAddress = &availabilitygrouplisteners.PrivateIPAddress{
			IPAddress:        pointer.To(lb.PrivateIpAddress),
			SubnetResourceId: pointer.To(lb.SubnetId),
		}

		results = append(results, lbConfig)
	}
	return &results, nil
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(multiSubnetIpConfiguration []MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener) *[]availabilitygrouplisteners.MultiSubnetIPConfiguration {
	results := make([]availabilitygrouplisteners.MultiSubnetIPConfiguration, 0)

	for _, item := range multiSubnetIpConfiguration {
		config := availabilitygrouplisteners.MultiSubnetIPConfiguration{
			SqlVirtualMachineInstance: item.SqlVirtualMachineId,
		}

		config.PrivateIPAddress = availabilitygrouplisteners.PrivateIPAddress{
			IPAddress:        pointer.To(item.PrivateIpAddress),
			SubnetResourceId: pointer.To(item.SubnetId),
		}

		results = append(results, config)
	}

	return &results
}

func expandMsSqlVirtualMachineAvailabilityGroupListenerReplicas(replicas []ReplicaMsSqlVirtualMachineAvailabilityGroupListener) (*[]availabilitygrouplisteners.AgReplica, error) {
	results := make([]availabilitygrouplisteners.AgReplica, 0)

	for _, rep := range replicas {
		replica := availabilitygrouplisteners.AgReplica{
			Role:              pointer.ToEnum[availabilitygrouplisteners.Role](rep.Role),
			Commit:            pointer.ToEnum[availabilitygrouplisteners.Commit](rep.Commit),
			Failover:          pointer.ToEnum[availabilitygrouplisteners.Failover](rep.FailoverMode),
			ReadableSecondary: pointer.ToEnum[availabilitygrouplisteners.ReadableSecondary](rep.ReadableSecondary),
		}

		sqlVirtualMachineId := rep.SqlVirtualMachineId
		if sqlVirtualMachineId != "" {
			id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(sqlVirtualMachineId)
			if err != nil {
				return nil, err
			}
			sqlVirtualMachineId = id.ID()
		}
		replica.SqlVirtualMachineInstanceId = pointer.To(sqlVirtualMachineId)

		results = append(results, replica)
	}
	return &results, nil
}
