package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProfile struct {
	AuthorizationProfile   AuthorizationProfile                  `json:"authorizationProfile"`
	AutoscaleProfile       *AutoscaleProfile                     `json:"autoscaleProfile,omitempty"`
	ClusterAccessProfile   *ClusterAccessProfile                 `json:"clusterAccessProfile,omitempty"`
	ClusterVersion         string                                `json:"clusterVersion"`
	Components             *[]ClusterComponentsComponentsInlined `json:"components,omitempty"`
	ConnectivityProfile    *ConnectivityProfile                  `json:"connectivityProfile,omitempty"`
	FlinkProfile           *FlinkProfile                         `json:"flinkProfile,omitempty"`
	IdentityProfile        *IdentityProfile                      `json:"identityProfile,omitempty"`
	KafkaProfile           *KafkaProfile                         `json:"kafkaProfile,omitempty"`
	LlapProfile            *interface{}                          `json:"llapProfile,omitempty"`
	LogAnalyticsProfile    *ClusterLogAnalyticsProfile           `json:"logAnalyticsProfile,omitempty"`
	ManagedIdentityProfile *ManagedIdentityProfile               `json:"managedIdentityProfile,omitempty"`
	OssVersion             string                                `json:"ossVersion"`
	PrometheusProfile      *ClusterPrometheusProfile             `json:"prometheusProfile,omitempty"`
	RangerPluginProfile    *ClusterRangerPluginProfile           `json:"rangerPluginProfile,omitempty"`
	RangerProfile          *RangerProfile                        `json:"rangerProfile,omitempty"`
	ScriptActionProfiles   *[]ScriptActionProfile                `json:"scriptActionProfiles,omitempty"`
	SecretsProfile         *SecretsProfile                       `json:"secretsProfile,omitempty"`
	ServiceConfigsProfiles *[]ClusterServiceConfigsProfile       `json:"serviceConfigsProfiles,omitempty"`
	SparkProfile           *SparkProfile                         `json:"sparkProfile,omitempty"`
	SshProfile             *SshProfile                           `json:"sshProfile,omitempty"`
	StubProfile            *interface{}                          `json:"stubProfile,omitempty"`
	TrinoProfile           *TrinoProfile                         `json:"trinoProfile,omitempty"`
}
