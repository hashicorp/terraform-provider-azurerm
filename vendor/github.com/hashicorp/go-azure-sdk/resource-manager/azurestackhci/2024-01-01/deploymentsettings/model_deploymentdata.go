package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentData struct {
	AdouPath              *string                     `json:"adouPath,omitempty"`
	Cluster               *DeploymentCluster          `json:"cluster,omitempty"`
	DomainFqdn            *string                     `json:"domainFqdn,omitempty"`
	HostNetwork           *HostNetwork                `json:"hostNetwork,omitempty"`
	InfrastructureNetwork *[]InfrastructureNetwork    `json:"infrastructureNetwork,omitempty"`
	NamingPrefix          *string                     `json:"namingPrefix,omitempty"`
	Observability         *Observability              `json:"observability,omitempty"`
	OptionalServices      *OptionalServices           `json:"optionalServices,omitempty"`
	PhysicalNodes         *[]PhysicalNodes            `json:"physicalNodes,omitempty"`
	SecretsLocation       *string                     `json:"secretsLocation,omitempty"`
	SecuritySettings      *DeploymentSecuritySettings `json:"securitySettings,omitempty"`
	Storage               *Storage                    `json:"storage,omitempty"`
}
