package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdatableClusterProfile struct {
	AuthorizationProfile   *AuthorizationProfile           `json:"authorizationProfile,omitempty"`
	AutoscaleProfile       *AutoscaleProfile               `json:"autoscaleProfile,omitempty"`
	LogAnalyticsProfile    *ClusterLogAnalyticsProfile     `json:"logAnalyticsProfile,omitempty"`
	PrometheusProfile      *ClusterPrometheusProfile       `json:"prometheusProfile,omitempty"`
	RangerPluginProfile    *ClusterRangerPluginProfile     `json:"rangerPluginProfile,omitempty"`
	RangerProfile          *RangerProfile                  `json:"rangerProfile,omitempty"`
	ScriptActionProfiles   *[]ScriptActionProfile          `json:"scriptActionProfiles,omitempty"`
	SecretsProfile         *SecretsProfile                 `json:"secretsProfile,omitempty"`
	ServiceConfigsProfiles *[]ClusterServiceConfigsProfile `json:"serviceConfigsProfiles,omitempty"`
	SshProfile             *SshProfile                     `json:"sshProfile,omitempty"`
	TrinoProfile           *TrinoProfile                   `json:"trinoProfile,omitempty"`
}
