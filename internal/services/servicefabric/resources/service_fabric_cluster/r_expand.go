// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_cluster

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandServiceFabricClusterAddOnFeatures(input []any) *[]cluster.AddOnFeatures {
	output := make([]cluster.AddOnFeatures, 0)

	for _, v := range input {
		output = append(output, cluster.AddOnFeatures(v.(string)))
	}

	return &output
}

func expandServiceFabricClusterAzureActiveDirectory(input []any) *cluster.AzureActiveDirectory {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	tenantId := v["tenant_id"].(string)
	clusterApplication := v["cluster_application_id"].(string)
	clientApplication := v["client_application_id"].(string)

	config := cluster.AzureActiveDirectory{
		TenantId:           pointer.To(tenantId),
		ClusterApplication: pointer.To(clusterApplication),
		ClientApplication:  pointer.To(clientApplication),
	}
	return &config
}

func expandServiceFabricClusterCertificate(input []any) *cluster.CertificateDescription {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	thumbprint := v["thumbprint"].(string)

	result := cluster.CertificateDescription{
		Thumbprint:    thumbprint,
		X509StoreName: pointer.ToEnum[cluster.X509StoreName](v["x509_store_name"].(string)),
	}

	if thumb, ok := v["thumbprint_secondary"]; ok {
		result.ThumbprintSecondary = pointer.To(thumb.(string))
	}

	return &result
}

func expandServiceFabricClusterCertificateCommonNames(d *pluginsdk.ResourceData) *cluster.ServerCertificateCommonNames {
	i := d.Get("certificate_common_names").([]any)
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]any)

	commonNamesRaw := input["common_names"].(*pluginsdk.Set).List()
	commonNames := make([]cluster.ServerCertificateCommonName, 0)

	for _, commonName := range commonNamesRaw {
		commonNameDetails := commonName.(map[string]any)

		commonName := cluster.ServerCertificateCommonName{
			CertificateCommonName:       commonNameDetails["certificate_common_name"].(string),
			CertificateIssuerThumbprint: commonNameDetails["certificate_issuer_thumbprint"].(string),
		}

		commonNames = append(commonNames, commonName)
	}

	output := cluster.ServerCertificateCommonNames{
		CommonNames:   &commonNames,
		X509StoreName: pointer.ToEnum[cluster.X509StoreName](input["x509_store_name"].(string)),
	}

	return &output
}

func expandServiceFabricClusterReverseProxyCertificateCommonNames(d *pluginsdk.ResourceData) *cluster.ServerCertificateCommonNames {
	i := d.Get("reverse_proxy_certificate_common_names").([]any)
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]any)

	commonNamesRaw := input["common_names"].(*pluginsdk.Set).List()
	commonNames := make([]cluster.ServerCertificateCommonName, 0)

	for _, commonName := range commonNamesRaw {
		commonNameDetails := commonName.(map[string]any)

		commonName := cluster.ServerCertificateCommonName{
			CertificateCommonName:       commonNameDetails["certificate_common_name"].(string),
			CertificateIssuerThumbprint: commonNameDetails["certificate_issuer_thumbprint"].(string),
		}

		commonNames = append(commonNames, commonName)
	}

	output := cluster.ServerCertificateCommonNames{
		CommonNames:   &commonNames,
		X509StoreName: pointer.ToEnum[cluster.X509StoreName](input["x509_store_name"].(string)),
	}

	return &output
}

func expandServiceFabricClusterReverseProxyCertificate(input []any) *cluster.CertificateDescription {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	result := cluster.CertificateDescription{
		Thumbprint:    v["thumbprint"].(string),
		X509StoreName: pointer.ToEnum[cluster.X509StoreName](v["x509_store_name"].(string)),
	}

	if thumb, ok := v["thumbprint_secondary"]; ok {
		result.ThumbprintSecondary = pointer.To(thumb.(string))
	}

	return &result
}

func expandServiceFabricClusterClientCertificateThumbprints(input []any) *[]cluster.ClientCertificateThumbprint {
	results := make([]cluster.ClientCertificateThumbprint, 0)

	for _, v := range input {
		val := v.(map[string]any)

		result := cluster.ClientCertificateThumbprint{
			CertificateThumbprint: val["thumbprint"].(string),
			IsAdmin:               val["is_admin"].(bool),
		}
		results = append(results, result)
	}

	return &results
}

func expandServiceFabricClusterClientCertificateCommonNames(input []any) *[]cluster.ClientCertificateCommonName {
	results := make([]cluster.ClientCertificateCommonName, 0)

	for _, v := range input {
		val := v.(map[string]any)

		certificateCommonName := val["common_name"].(string)
		certificateIssuerThumbprint := val["issuer_thumbprint"].(string)
		isAdmin := val["is_admin"].(bool)

		result := cluster.ClientCertificateCommonName{
			CertificateCommonName:       certificateCommonName,
			CertificateIssuerThumbprint: certificateIssuerThumbprint,
			IsAdmin:                     isAdmin,
		}
		results = append(results, result)
	}

	return &results
}

func expandServiceFabricClusterDiagnosticsConfig(input []any) *cluster.DiagnosticsStorageAccountConfig {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	config := cluster.DiagnosticsStorageAccountConfig{
		StorageAccountName:      v["storage_account_name"].(string),
		ProtectedAccountKeyName: v["protected_account_key_name"].(string),
		BlobEndpoint:            v["blob_endpoint"].(string),
		QueueEndpoint:           v["queue_endpoint"].(string),
		TableEndpoint:           v["table_endpoint"].(string),
	}
	return &config
}

func expandServiceFabricClusterUpgradePolicyDeltaHealthPolicy(input []any) *cluster.ClusterUpgradeDeltaHealthPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	deltaHealthPolicy := &cluster.ClusterUpgradeDeltaHealthPolicy{}
	v := input[0].(map[string]any)
	deltaHealthPolicy.MaxPercentDeltaUnhealthyNodes = int64(v["max_delta_unhealthy_nodes_percent"].(int))
	deltaHealthPolicy.MaxPercentUpgradeDomainDeltaUnhealthyNodes = int64(v["max_upgrade_domain_delta_unhealthy_nodes_percent"].(int))
	deltaHealthPolicy.MaxPercentDeltaUnhealthyApplications = int64(v["max_delta_unhealthy_applications_percent"].(int))

	return deltaHealthPolicy
}

func expandServiceFabricClusterUpgradePolicyHealthPolicy(input []any) cluster.ClusterHealthPolicy {
	healthPolicy := cluster.ClusterHealthPolicy{}
	if len(input) == 0 || input[0] == nil {
		return healthPolicy
	}

	v := input[0].(map[string]any)
	healthPolicy.MaxPercentUnhealthyApplications = pointer.To(int64(v["max_unhealthy_applications_percent"].(int)))
	healthPolicy.MaxPercentUnhealthyNodes = pointer.To(int64(v["max_unhealthy_nodes_percent"].(int)))

	return healthPolicy
}

func expandServiceFabricClusterUpgradePolicy(input []any) *cluster.ClusterUpgradePolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	policy := &cluster.ClusterUpgradePolicy{}
	v := input[0].(map[string]any)

	policy.ForceRestart = pointer.To(v["force_restart_enabled"].(bool))
	policy.HealthCheckStableDuration = v["health_check_stable_duration"].(string)
	policy.UpgradeDomainTimeout = v["upgrade_domain_timeout"].(string)
	policy.UpgradeReplicaSetCheckTimeout = v["upgrade_replica_set_check_timeout"].(string)
	policy.UpgradeTimeout = v["upgrade_timeout"].(string)
	policy.HealthCheckRetryTimeout = v["health_check_retry_timeout"].(string)
	policy.HealthCheckWaitDuration = v["health_check_wait_duration"].(string)

	if v["health_policy"] != nil {
		policy.HealthPolicy = expandServiceFabricClusterUpgradePolicyHealthPolicy(v["health_policy"].([]any))
	}
	if v["delta_health_policy"] != nil {
		policy.DeltaHealthPolicy = expandServiceFabricClusterUpgradePolicyDeltaHealthPolicy(v["delta_health_policy"].([]any))
	}

	return policy
}

func expandServiceFabricClusterFabricSettings(input []any) *[]cluster.SettingsSectionDescription {
	results := make([]cluster.SettingsSectionDescription, 0)

	for _, v := range input {
		val := v.(map[string]any)

		name := val["name"].(string)
		params := make([]cluster.SettingsParameterDescription, 0)
		paramsRaw := val["parameters"].(map[string]any)
		for k, v := range paramsRaw {
			param := cluster.SettingsParameterDescription{
				Name:  k,
				Value: v.(string),
			}
			params = append(params, param)
		}

		result := cluster.SettingsSectionDescription{
			Name:       name,
			Parameters: params,
		}
		results = append(results, result)
	}

	return &results
}

func expandServiceFabricClusterNodeTypes(input []any) []cluster.NodeTypeDescription {
	results := make([]cluster.NodeTypeDescription, 0)

	for _, v := range input {
		node := v.(map[string]any)

		result := cluster.NodeTypeDescription{
			Name:                         node["name"].(string),
			VMInstanceCount:              int64(node["instance_count"].(int)),
			IsPrimary:                    node["is_primary"].(bool),
			ClientConnectionEndpointPort: int64(node["client_endpoint_port"].(int)),
			HTTPGatewayEndpointPort:      int64(node["http_endpoint_port"].(int)),
			DurabilityLevel:              pointer.ToEnum[cluster.DurabilityLevel](node["durability_level"].(string)),
		}

		if isStateless, ok := node["is_stateless"]; ok {
			result.IsStateless = pointer.To(isStateless.(bool))
		}

		if multipleAvailabilityZones, ok := node["multiple_availability_zones"]; ok {
			result.MultipleAvailabilityZones = pointer.To(multipleAvailabilityZones.(bool))
		}

		if props, ok := node["placement_properties"]; ok {
			placementProperties := make(map[string]string)
			for key, value := range props.(map[string]any) {
				placementProperties[key] = value.(string)
			}

			result.PlacementProperties = &placementProperties
		}

		if caps, ok := node["capacities"]; ok {
			capacities := make(map[string]string)
			for key, value := range caps.(map[string]any) {
				capacities[key] = value.(string)
			}

			result.Capacities = &capacities
		}

		if v := int64(node["reverse_proxy_endpoint_port"].(int)); v != 0 {
			result.ReverseProxyEndpointPort = pointer.To(v)
		}

		applicationPortsRaw := node["application_ports"].([]any)
		if len(applicationPortsRaw) > 0 {
			portsRaw := applicationPortsRaw[0].(map[string]any)

			result.ApplicationPorts = &cluster.EndpointRangeDescription{
				StartPort: int64(portsRaw["start_port"].(int)),
				EndPort:   int64(portsRaw["end_port"].(int)),
			}
		}

		ephemeralPortsRaw := node["ephemeral_ports"].([]any)
		if len(ephemeralPortsRaw) > 0 {
			portsRaw := ephemeralPortsRaw[0].(map[string]any)

			result.EphemeralPorts = &cluster.EndpointRangeDescription{
				StartPort: int64(portsRaw["start_port"].(int)),
				EndPort:   int64(portsRaw["end_port"].(int)),
			}
		}

		results = append(results, result)
	}

	return results
}
