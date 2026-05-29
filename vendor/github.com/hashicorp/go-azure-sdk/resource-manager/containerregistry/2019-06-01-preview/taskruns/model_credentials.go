package taskruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Credentials struct {
	CustomRegistries *map[string]CustomRegistryCredentials `json:"customRegistries,omitempty"`
	SourceRegistry   *SourceRegistryCredentials            `json:"sourceRegistry,omitempty"`
}
