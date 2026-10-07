package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasPathMetadata struct {
	Attributes *AliasPathAttributes `json:"attributes,omitempty"`
	Type       *AliasPathTokenType  `json:"type,omitempty"`
}
