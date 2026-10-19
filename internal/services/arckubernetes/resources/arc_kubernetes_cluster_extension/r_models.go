// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster_extension

type ArcKubernetesClusterExtensionModel struct {
	Name                           string            `tfschema:"name"`
	ClusterID                      string            `tfschema:"cluster_id"`
	ConfigurationProtectedSettings map[string]string `tfschema:"configuration_protected_settings"`
	ConfigurationSettings          map[string]string `tfschema:"configuration_settings"`
	ExtensionType                  string            `tfschema:"extension_type"`
	ReleaseNamespace               string            `tfschema:"release_namespace"`
	ReleaseTrain                   string            `tfschema:"release_train"`
	TargetNamespace                string            `tfschema:"target_namespace"`
	Version                        string            `tfschema:"version"`
	CurrentVersion                 string            `tfschema:"current_version"`
}

func (r Resource) ModelObject() any {
	return &ArcKubernetesClusterExtensionModel{}
}
