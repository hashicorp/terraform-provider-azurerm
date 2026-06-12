package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Configuration struct {
	ActiveRevisionsMode  *ActiveRevisionsMode   `json:"activeRevisionsMode,omitempty"`
	Dapr                 *Dapr                  `json:"dapr,omitempty"`
	IdentitySettings     *[]IdentitySettings    `json:"identitySettings,omitempty"`
	Ingress              *Ingress               `json:"ingress,omitempty"`
	MaxInactiveRevisions *int64                 `json:"maxInactiveRevisions,omitempty"`
	Registries           *[]RegistryCredentials `json:"registries,omitempty"`
	Runtime              *Runtime               `json:"runtime,omitempty"`
	Secrets              *[]Secret              `json:"secrets,omitempty"`
	Service              *Service               `json:"service,omitempty"`
}
