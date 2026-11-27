package containerappssessionpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomContainerTemplate struct {
	Containers          *[]SessionContainer         `json:"containers,omitempty"`
	Ingress             *SessionIngress             `json:"ingress,omitempty"`
	RegistryCredentials *SessionRegistryCredentials `json:"registryCredentials,omitempty"`
}
