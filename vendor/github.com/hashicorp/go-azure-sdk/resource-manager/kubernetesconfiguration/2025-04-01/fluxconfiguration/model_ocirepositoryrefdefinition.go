package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OCIRepositoryRefDefinition struct {
	Digest *string `json:"digest,omitempty"`
	Semver *string `json:"semver,omitempty"`
	Tag    *string `json:"tag,omitempty"`
}
