package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RepositoryRefDefinition struct {
	Branch *string `json:"branch,omitempty"`
	Commit *string `json:"commit,omitempty"`
	Semver *string `json:"semver,omitempty"`
	Tag    *string `json:"tag,omitempty"`
}
