// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_availability_group_listener

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
)

func flattenMsSqlVirtualMachineAvailabilityGroupListenerLoadBalancerConfigurations(input *[]availabilitygrouplisteners.LoadBalancerConfiguration, subscriptionId string) ([]LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, lbConfig := range *input {
		privateIpAddress := ""
		subnetResourceId := ""
		if v := lbConfig.PrivateIPAddress; v != nil {
			privateIpAddress = pointer.From(v.IPAddress)

			parsedSubnetResourceId, err := parse.SubnetIDInsensitively(pointer.From(v.SubnetResourceId))
			if err != nil {
				return nil, err
			}
			subnetResourceId = parsedSubnetResourceId.ID()
		}

		loadBalancerId := ""
		if lbConfig.LoadBalancerResourceId != nil {
			id, err := loadbalancers.ParseLoadBalancerIDInsensitively(pointer.From(lbConfig.LoadBalancerResourceId))
			if err != nil {
				return nil, err
			}
			loadBalancerId = id.ID()
		}

		var sqlVirtualMachineIds []string
		if lbConfig.SqlVirtualMachineInstances != nil {
			sqlVirtualMachineIds = *lbConfig.SqlVirtualMachineInstances
			var parsedIds []string
			for _, sqlVmId := range sqlVirtualMachineIds {
				parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively(sqlVmId)
				if err != nil {
					return nil, err
				}

				// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
				newId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineName)

				parsedIds = append(parsedIds, newId.ID())
			}
			sqlVirtualMachineIds = parsedIds
		}

		v := LoadBalancerConfigurationMsSqlVirtualMachineAvailabilityGroupListener{
			LoadBalancerId:       loadBalancerId,
			PrivateIpAddress:     privateIpAddress,
			ProbePort:            pointer.From(lbConfig.ProbePort),
			SqlVirtualMachineIds: sqlVirtualMachineIds,
			SubnetId:             subnetResourceId,
		}

		results = append(results, v)
	}
	return results, nil
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerMultiSubnetIpConfiguration(input *[]availabilitygrouplisteners.MultiSubnetIPConfiguration, subscriptionId string) ([]MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, config := range *input {
		parsedSubnetResourceId, err := parse.SubnetIDInsensitively(pointer.From(config.PrivateIPAddress.SubnetResourceId))
		if err != nil {
			return nil, err
		}

		parsedSqlVirtualMachineId, err := sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively(config.SqlVirtualMachineInstance)
		if err != nil {
			return nil, err
		}

		// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
		newSqlVirtualMachineId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedSqlVirtualMachineId.ResourceGroupName, parsedSqlVirtualMachineId.SqlVirtualMachineName)

		v := MultiSubnetIpConfigurationMsSqlVirtualMachineAvailabilityGroupListener{
			PrivateIpAddress:    pointer.From(config.PrivateIPAddress.IPAddress),
			SqlVirtualMachineId: newSqlVirtualMachineId.ID(),
			SubnetId:            parsedSubnetResourceId.ID(),
		}

		results = append(results, v)
	}
	return results, nil
}

func flattenMsSqlVirtualMachineAvailabilityGroupListenerReplicas(input *[]availabilitygrouplisteners.AgReplica, subscriptionId string) ([]ReplicaMsSqlVirtualMachineAvailabilityGroupListener, error) {
	results := make([]ReplicaMsSqlVirtualMachineAvailabilityGroupListener, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, replica := range *input {
		sqlVirtualMachineInstanceId := ""
		if replica.SqlVirtualMachineInstanceId != nil {
			parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively(*replica.SqlVirtualMachineInstanceId)
			if err != nil {
				return nil, err
			}

			// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
			newId := sqlvirtualmachines.NewSqlVirtualMachineID(subscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineName)

			sqlVirtualMachineInstanceId = newId.ID()
		}

		v := ReplicaMsSqlVirtualMachineAvailabilityGroupListener{
			SqlVirtualMachineId: sqlVirtualMachineInstanceId,
			Role:                pointer.FromEnum(replica.Role),
			Commit:              pointer.FromEnum(replica.Commit),
			FailoverMode:        pointer.FromEnum(replica.Failover),
			ReadableSecondary:   pointer.FromEnum(replica.ReadableSecondary),
		}

		results = append(results, v)
	}
	return results, nil
}
