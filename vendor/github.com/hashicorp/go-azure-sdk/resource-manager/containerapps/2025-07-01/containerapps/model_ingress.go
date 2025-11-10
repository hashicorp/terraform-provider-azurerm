package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Ingress struct {
	AdditionalPortMappings *[]IngressPortMapping         `json:"additionalPortMappings,omitempty"`
	AllowInsecure          *bool                         `json:"allowInsecure,omitempty"`
	ClientCertificateMode  *IngressClientCertificateMode `json:"clientCertificateMode,omitempty"`
	CorsPolicy             *CorsPolicy                   `json:"corsPolicy,omitempty"`
	CustomDomains          *[]CustomDomain               `json:"customDomains,omitempty"`
	ExposedPort            *int64                        `json:"exposedPort,omitempty"`
	External               *bool                         `json:"external,omitempty"`
	Fqdn                   *string                       `json:"fqdn,omitempty"`
	IPSecurityRestrictions *[]IPSecurityRestrictionRule  `json:"ipSecurityRestrictions,omitempty"`
	StickySessions         *IngressStickySessions        `json:"stickySessions,omitempty"`
	TargetPort             *int64                        `json:"targetPort,omitempty"`
	Traffic                *[]TrafficWeight              `json:"traffic,omitempty"`
	Transport              *IngressTransportMethod       `json:"transport,omitempty"`
}
