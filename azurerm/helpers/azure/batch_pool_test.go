package azure

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandBatchPoolNetworkConfiguration(t *testing.T) {

	networkSecurityGroupRule := make(map[string]interface{})
	networkSecurityGroupRule["priority"] = 150
	networkSecurityGroupRule["access"] = "Allow"
	networkSecurityGroupRule["source_address_prefix"] = "prefix"

	networkSecurityGroupRules := make([]interface{}, 1)
	networkSecurityGroupRules[0] = networkSecurityGroupRule

	inboundNatPool := make(map[string]interface{})
	inboundNatPool["name"] = "Name"
	inboundNatPool["protocol"] = "TCP"
	inboundNatPool["backend_port"] = 2
	inboundNatPool["frontend_port_range_start"] = 3
	inboundNatPool["frontend_port_range_end"] = 6
	inboundNatPool["network_security_group_rules"] = networkSecurityGroupRules

	inboundNatPools := make([]interface{}, 1)
	inboundNatPools[0] = inboundNatPool

	endpointConfig := make(map[string]interface{})
	endpointConfig["inbound_nat_pools"] = inboundNatPools

	networkConfig := make(map[string]interface{})
	networkConfig["subnet_id"] = "test"
	networkConfig["endpoint_configuration"] = make([]interface{}, 1)
	networkConfig["endpoint_configuration"].([]interface{})[0] = endpointConfig

	input := make([]interface{}, 1)
	input[0] = networkConfig

	cases := []struct {
		Input       []interface{}
		ExpectError bool
	}{
		{
			Input:       input,
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		nc, err := ExpandBatchPoolNetworkConfiguration(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}
			return
		}

		assert.EqualValues(t, "test", *nc.SubnetID)
		assert.EqualValues(t, 1, len(*nc.EndpointConfiguration.InboundNatPools))

		inboundNatPools := (*nc.EndpointConfiguration.InboundNatPools)[0]
		assert.EqualValues(t, batch.TCP, inboundNatPools.Protocol)
		assert.EqualValues(t, "Name", *inboundNatPools.Name)
		assert.EqualValues(t, 3, *inboundNatPools.FrontendPortRangeStart)
		assert.EqualValues(t, 6, *inboundNatPools.FrontendPortRangeEnd)
		assert.EqualValues(t, 2, *inboundNatPools.BackendPort)

		assert.Equal(t, 1, len(*inboundNatPools.NetworkSecurityGroupRules))

		groupRules := (*inboundNatPools.NetworkSecurityGroupRules)[0]

		assert.EqualValues(t, batch.Allow, groupRules.Access)
		assert.EqualValues(t, 150, *groupRules.Priority)
		assert.EqualValues(t, "prefix", *groupRules.SourceAddressPrefix)

	}
}

func TestExpandBatchPoolNetworkConfigurationOnlySubnetId(t *testing.T) {

	var jsonBlob = []byte(`[
  {
    "subnet_id": "test"
  }
]`)
	var input []interface{}
	err := json.Unmarshal(jsonBlob, &input)
	if err != nil {
		t.Fatalf("Got error when unmarshaling %+v", err)
	}

	cases := []struct {
		Input       []interface{}
		ExpectError bool
	}{
		{
			Input:       input,
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		nc, err := ExpandBatchPoolNetworkConfiguration(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}

			return
		}

		assert.EqualValues(t, "test", *nc.SubnetID)
		assert.True(t, nil == nc.EndpointConfiguration)
	}
}

func TestFlattenBatchPoolNetworkConfiguration(t *testing.T) {

	subnetId := "subnetId"
	backendPort := int32(1)
	frontendPortRangeStart := int32(2)
	frontendPortRangeEnd := int32(3)

	prefix := "prefix"
	priority := int32(99)
	groupRules := make([]batch.NetworkSecurityGroupRule, 1)
	groupRules[0].Access = batch.Allow
	groupRules[0].SourceAddressPrefix = &prefix
	groupRules[0].Priority = &priority

	inboundNatPool := make([]batch.InboundNatPool, 1)
	inboundNatPool[0].BackendPort = &backendPort
	inboundNatPool[0].FrontendPortRangeStart = &frontendPortRangeStart
	inboundNatPool[0].FrontendPortRangeEnd = &frontendPortRangeEnd
	inboundNatPool[0].Protocol = batch.UDP
	inboundNatPool[0].NetworkSecurityGroupRules = &groupRules

	networkConfiguration := &batch.NetworkConfiguration{
		SubnetID: &subnetId,
		EndpointConfiguration: &batch.PoolEndpointConfiguration{
			InboundNatPools: &inboundNatPool,
		},
	}

	flatten := FlattenBatchPoolNetworkConfiguration(networkConfiguration)

	assert.EqualValues(t, "[map[endpoint_configuration:[map[inbound_nat_pools:[map[backend_port:1 frontend_port_range_end:3 frontend_port_range_start:2 network_security_group_rules:[map[access:Allow priority:99 source_address_prefix:prefix]] protocol:UDP]]]] subnet_id:subnetId]]",
		fmt.Sprintf("%v", flatten))
}

func TestFlattenBatchPoolNetworkConfigurationEmpty(t *testing.T) {

	inboundNatPool := make([]batch.InboundNatPool, 1)
	groupRules := make([]batch.NetworkSecurityGroupRule, 1)

	inboundNatPool[0].NetworkSecurityGroupRules = &groupRules

	networkConfiguration := &batch.NetworkConfiguration{
		EndpointConfiguration: &batch.PoolEndpointConfiguration{
			InboundNatPools: &inboundNatPool,
		},
	}

	flatten := FlattenBatchPoolNetworkConfiguration(networkConfiguration)

	assert.EqualValues(t, "[map[endpoint_configuration:[map[inbound_nat_pools:[map[network_security_group_rules:[map[access:]] protocol:]]]]]]",
		fmt.Sprintf("%v", flatten))
}
