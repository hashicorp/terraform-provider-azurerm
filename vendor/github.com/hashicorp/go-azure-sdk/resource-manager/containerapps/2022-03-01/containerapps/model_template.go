package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Template struct {
	Containers     *[]Container `json:"containers,omitempty"`
	RevisionSuffix *string      `json:"revisionSuffix,omitempty"`
	Scale          *Scale       `json:"scale,omitempty"`
	Volumes        *[]Volume    `json:"volumes,omitempty"`
}
