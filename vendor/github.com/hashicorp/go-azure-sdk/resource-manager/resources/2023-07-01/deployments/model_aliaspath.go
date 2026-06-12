package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasPath struct {
	ApiVersions *[]string          `json:"apiVersions,omitempty"`
	Metadata    *AliasPathMetadata `json:"metadata,omitempty"`
	Path        *string            `json:"path,omitempty"`
	Pattern     *AliasPattern      `json:"pattern,omitempty"`
}
