// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_extension

import (
	"fmt"
	"maps"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2024-11-01/extensions"
)

func kubernetesClusterExtensionSettingsRemoved(previous map[string]any, settings map[string]string) bool {
	for key := range previous {
		if _, ok := settings[key]; !ok {
			return true
		}
	}
	return false
}

func expandKubernetesClusterExtensionFullUpdate(existing extensions.Extension, settings, protectedSettings map[string]string, previousProtectedSettings map[string]any) (extensions.Extension, error) {
	properties := existing.Properties
	if properties == nil || properties.ExtensionType == nil || properties.AutoUpgradeMinorVersion == nil || properties.Scope == nil {
		return extensions.Extension{}, fmt.Errorf("cannot replace extension configuration: missing existing extension properties")
	}
	if properties.ConfigurationProtectedSettings == nil {
		return extensions.Extension{}, fmt.Errorf("cannot replace extension configuration: the API did not return configurationProtectedSettings, so protected settings cannot be verified")
	}

	// Azure returns protected setting names, but not their values. Never replay
	// the redacted response or implicitly discard settings unknown to Terraform.
	for key := range *properties.ConfigurationProtectedSettings {
		if _, ok := protectedSettings[key]; ok {
			continue
		}
		if _, ok := previousProtectedSettings[key]; !ok {
			return extensions.Extension{}, fmt.Errorf("cannot replace extension configuration: existing protected settings are not available in Terraform state; configure configuration_protected_settings before removing configuration_settings")
		}
	}

	// Use explicit empty objects when removing the last setting, rather than
	// sending null maps whose replacement semantics may differ.
	configuration := make(map[string]string, len(settings))
	maps.Copy(configuration, settings)
	protected := make(map[string]string, len(protectedSettings))
	maps.Copy(protected, protectedSettings)

	// Build a writable payload, not a copy of GET. In particular, generated
	// principal IDs, provisioning status, and system metadata must not be sent.
	result := extensions.Extension{
		Identity: existing.Identity,
		Plan:     existing.Plan,
		Properties: &extensions.ExtensionProperties{
			ExtensionType:                  properties.ExtensionType,
			AutoUpgradeMinorVersion:        properties.AutoUpgradeMinorVersion,
			ReleaseTrain:                   properties.ReleaseTrain,
			Version:                        properties.Version,
			Scope:                          properties.Scope,
			ConfigurationSettings:          &configuration,
			ConfigurationProtectedSettings: &protected,
		},
	}
	// SystemAssigned's JSON marshaler only emits the writable type, never the
	// generated IDs. AKS identity does not have that custom marshaler.
	if properties.AksAssignedIdentity != nil && properties.AksAssignedIdentity.Type != nil {
		result.Properties.AksAssignedIdentity = &extensions.ExtensionPropertiesAksAssignedIdentity{
			Type: properties.AksAssignedIdentity.Type,
		}
	}
	return result, nil
}
