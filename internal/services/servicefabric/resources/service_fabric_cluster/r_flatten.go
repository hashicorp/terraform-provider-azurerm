// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_cluster

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func flattenServiceFabricClusterAzureActiveDirectory(input *cluster.AzureActiveDirectory) []any {
	results := make([]any, 0)

	if v := input; v != nil {
		output := make(map[string]any)

		if name := v.TenantId; name != nil {
			output["tenant_id"] = *name
		}

		if name := v.ClusterApplication; name != nil {
			output["cluster_application_id"] = *name
		}

		if endpoint := v.ClientApplication; endpoint != nil {
			output["client_application_id"] = *endpoint
		}

		results = append(results, output)
	}

	return results
}

func flattenServiceFabricClusterAddOnFeatures(input *[]cluster.AddOnFeatures) []any {
	output := make([]any, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, string(v))
		}
	}

	return output
}

func flattenServiceFabricClusterCertificate(input *cluster.CertificateDescription) []any {
	results := make([]any, 0)

	if v := input; v != nil {
		output := make(map[string]any)

		output["thumbprint"] = input.Thumbprint

		if thumbprint := input.ThumbprintSecondary; thumbprint != nil {
			output["thumbprint_secondary"] = *thumbprint
		}

		output["x509_store_name"] = input.X509StoreName
		results = append(results, output)
	}

	return results
}

func flattenServiceFabricClusterCertificateCommonNames(in *cluster.ServerCertificateCommonNames) []any {
	if in == nil {
		return []any{}
	}

	output := make(map[string]any)

	if commonNames := in.CommonNames; commonNames != nil {
		common_names := make([]map[string]any, 0)
		for _, i := range *commonNames {
			commonName := make(map[string]any)

			commonName["certificate_common_name"] = i.CertificateCommonName
			commonName["certificate_issuer_thumbprint"] = i.CertificateIssuerThumbprint

			common_names = append(common_names, commonName)
		}

		output["common_names"] = common_names
	}

	output["x509_store_name"] = in.X509StoreName

	return []any{output}
}

func flattenServiceFabricClusterReverseProxyCertificate(input *cluster.CertificateDescription) []any {
	results := make([]any, 0)

	if v := input; v != nil {
		output := make(map[string]any)

		output["thumbprint"] = input.Thumbprint

		if thumbprint := input.ThumbprintSecondary; thumbprint != nil {
			output["thumbprint_secondary"] = *thumbprint
		}

		output["x509_store_name"] = input.X509StoreName
		results = append(results, output)
	}

	return results
}

func flattenServiceFabricClusterClientCertificateThumbprints(input *[]cluster.ClientCertificateThumbprint) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)

	for _, v := range *input {
		result := make(map[string]any)

		result["thumbprint"] = v.CertificateThumbprint
		result["admin"] = v.IsAdmin
		if !features.SixPointOh() {
			result["is_admin"] = v.IsAdmin
		}

		results = append(results, result)
	}

	return results
}

func flattenServiceFabricClusterClientCertificateCommonNames(input *[]cluster.ClientCertificateCommonName) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)

	for _, v := range *input {
		result := make(map[string]any)

		result["common_name"] = v.CertificateCommonName
		result["issuer_thumbprint"] = v.CertificateIssuerThumbprint
		result["admin"] = v.IsAdmin
		if !features.SixPointOh() {
			result["is_admin"] = v.IsAdmin
		}

		results = append(results, result)
	}

	return results
}

func flattenServiceFabricClusterDiagnosticsConfig(input *cluster.DiagnosticsStorageAccountConfig) []any {
	results := make([]any, 0)

	if v := input; v != nil {
		output := make(map[string]any)

		output["storage_account_name"] = v.StorageAccountName
		output["protected_account_key_name"] = v.ProtectedAccountKeyName
		output["blob_endpoint"] = v.BlobEndpoint
		output["queue_endpoint"] = v.QueueEndpoint
		output["table_endpoint"] = v.TableEndpoint

		results = append(results, output)
	}

	return results
}

func flattenServiceFabricClusterUpgradePolicy(input *cluster.ClusterUpgradePolicy) []any {
	if input == nil {
		return []any{}
	}

	output := make(map[string]any)

	if forceRestart := input.ForceRestart; forceRestart != nil {
		output["force_restart_enabled"] = *forceRestart
	}

	output["health_check_retry_timeout"] = input.HealthCheckRetryTimeout
	output["health_check_stable_duration"] = input.HealthCheckStableDuration
	output["health_check_wait_duration"] = input.HealthCheckWaitDuration
	output["upgrade_domain_timeout"] = input.UpgradeDomainTimeout
	output["upgrade_replica_set_check_timeout"] = input.UpgradeReplicaSetCheckTimeout
	output["upgrade_timeout"] = input.UpgradeTimeout

	output["health_policy"] = flattenServiceFabricClusterUpgradePolicyHealthPolicy(input.HealthPolicy)
	output["delta_health_policy"] = flattenServiceFabricClusterUpgradePolicyDeltaHealthPolicy(input.DeltaHealthPolicy)

	return []any{output}
}

func flattenServiceFabricClusterUpgradePolicyHealthPolicy(input cluster.ClusterHealthPolicy) []any {
	output := make(map[string]any)

	if input.MaxPercentUnhealthyApplications != nil {
		output["max_unhealthy_applications_percent"] = *input.MaxPercentUnhealthyApplications
	}

	if input.MaxPercentUnhealthyNodes != nil {
		output["max_unhealthy_nodes_percent"] = *input.MaxPercentUnhealthyNodes
	}

	return []any{output}
}

func flattenServiceFabricClusterUpgradePolicyDeltaHealthPolicy(input *cluster.ClusterUpgradeDeltaHealthPolicy) []any {
	if input == nil {
		return []any{}
	}

	output := make(map[string]any)

	output["max_delta_unhealthy_applications_percent"] = input.MaxPercentDeltaUnhealthyApplications
	output["max_delta_unhealthy_nodes_percent"] = input.MaxPercentDeltaUnhealthyNodes
	output["max_upgrade_domain_delta_unhealthy_nodes_percent"] = input.MaxPercentUpgradeDomainDeltaUnhealthyNodes

	return []any{output}
}

func flattenServiceFabricClusterFabricSettings(input *[]cluster.SettingsSectionDescription) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)

	for _, v := range *input {
		result := make(map[string]any)

		result["name"] = v.Name

		parameters := make(map[string]any)
		if paramsRaw := v.Parameters; paramsRaw != nil {
			for _, p := range paramsRaw {
				parameters[p.Name] = p.Value
			}
		}
		result["parameters"] = parameters
		results = append(results, result)
	}

	return results
}

func flattenServiceFabricClusterNodeTypes(input []cluster.NodeTypeDescription) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)

	for _, v := range input {
		output := make(map[string]any)

		output["name"] = v.Name
		output["instance_count"] = v.VMInstanceCount
		output["primary"] = v.IsPrimary
		if !features.SixPointOh() {
			output["is_primary"] = v.IsPrimary
		}

		output["client_endpoint_port"] = v.ClientConnectionEndpointPort
		output["http_endpoint_port"] = v.HTTPGatewayEndpointPort

		if placementProperties := v.PlacementProperties; placementProperties != nil {
			output["placement_properties"] = *placementProperties
		}

		if capacities := v.Capacities; capacities != nil {
			output["capacities"] = *capacities
		}

		if port := v.ReverseProxyEndpointPort; port != nil {
			output["reverse_proxy_endpoint_port"] = *port
		}

		if isStateless := v.IsStateless; isStateless != nil {
			output["stateless"] = *isStateless
			if !features.SixPointOh() {
				output["is_stateless"] = *isStateless
			}
		}

		if multipleAvailabilityZones := v.MultipleAvailabilityZones; multipleAvailabilityZones != nil {
			output["multiple_availability_zones"] = *multipleAvailabilityZones
		}

		if durabilityLevel := v.DurabilityLevel; durabilityLevel != nil {
			output["durability_level"] = string(*v.DurabilityLevel)
		}

		applicationPorts := make([]any, 0)
		if ports := v.ApplicationPorts; ports != nil {
			r := make(map[string]any)
			r["start_port"] = int(ports.StartPort)
			r["end_port"] = int(ports.EndPort)
			applicationPorts = append(applicationPorts, r)
		}
		output["application_ports"] = applicationPorts

		ephemeralPorts := make([]any, 0)
		if ports := v.EphemeralPorts; ports != nil {
			r := make(map[string]any)
			r["start_port"] = int(ports.StartPort)
			r["end_port"] = int(ports.EndPort)
			ephemeralPorts = append(ephemeralPorts, r)
		}
		output["ephemeral_ports"] = ephemeralPorts

		results = append(results, output)
	}

	return results
}
