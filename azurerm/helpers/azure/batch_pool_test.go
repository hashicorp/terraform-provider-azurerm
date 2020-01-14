package azure

import (
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
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
		_, err := ExpandBatchPoolNetworkConfiguration(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}
			return
		}
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
		_, err := ExpandBatchPoolNetworkConfiguration(tc.Input)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}
			return
		}
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

	FlattenBatchPoolNetworkConfiguration(networkConfiguration)
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

	FlattenBatchPoolNetworkConfiguration(networkConfiguration)
}
