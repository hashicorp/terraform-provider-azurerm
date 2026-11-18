// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
)

func ExpandHDInsightApplicationEdgeNodeInstallScriptActions(input []any) *[]applications.RuntimeScriptAction {
	actions := make([]applications.RuntimeScriptAction, 0)

	for _, v := range input {
		val := v.(map[string]any)

		name := val["name"].(string)
		uri := val["uri"].(string)
		parameters := val["parameters"].(string)

		action := applications.RuntimeScriptAction{
			Name: name,
			Uri:  uri,
			// The only role available for edge nodes is edgenode
			Parameters: pointer.To(parameters),
			Roles:      []string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}

func ExpandHDInsightApplicationEdgeNodeHttpsEndpoints(input []any) *[]applications.ApplicationGetHTTPSEndpoint {
	endpoints := make([]applications.ApplicationGetHTTPSEndpoint, 0)
	if len(input) == 0 || input[0] == nil {
		return &endpoints
	}

	for _, v := range input {
		val := v.(map[string]any)

		destinationPort := val["destination_port"].(int64)
		disableGatewayAuth := val["disable_gateway_auth"].(bool)
		privateIpAddress := val["private_ip_address"].(string)
		subDomainSuffix := val["sub_domain_suffix"].(string)

		endPoint := applications.ApplicationGetHTTPSEndpoint{
			AccessModes:        pointer.To(val["access_modes"].([]string)),
			DestinationPort:    pointer.To(destinationPort),
			PrivateIPAddress:   pointer.To(privateIpAddress),
			SubDomainSuffix:    pointer.To(subDomainSuffix),
			DisableGatewayAuth: pointer.To(disableGatewayAuth),
		}

		endpoints = append(endpoints, endPoint)
	}

	return &endpoints
}

func ExpandHDInsightApplicationEdgeNodeUninstallScriptActions(input []any) *[]applications.RuntimeScriptAction {
	actions := make([]applications.RuntimeScriptAction, 0)
	if len(input) == 0 || input[0] == nil {
		return &actions
	}

	for _, v := range input {
		val := v.(map[string]any)

		name := val["name"].(string)
		uri := val["uri"].(string)
		parameters := val["parameters"].(string)

		action := applications.RuntimeScriptAction{
			Name:       name,
			Uri:        uri,
			Parameters: pointer.To(parameters),
			Roles:      []string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}
