package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Ingress struct {
	AllowInsecure *bool                   `json:"allowInsecure,omitempty"`
	CustomDomains *[]CustomDomain         `json:"customDomains,omitempty"`
	External      *bool                   `json:"external,omitempty"`
	Fqdn          *string                 `json:"fqdn,omitempty"`
	TargetPort    *int64                  `json:"targetPort,omitempty"`
	Traffic       *[]TrafficWeight        `json:"traffic,omitempty"`
	Transport     *IngressTransportMethod `json:"transport,omitempty"`
}
