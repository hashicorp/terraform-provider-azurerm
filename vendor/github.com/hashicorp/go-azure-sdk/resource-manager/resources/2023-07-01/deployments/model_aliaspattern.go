package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasPattern struct {
	Phrase   *string           `json:"phrase,omitempty"`
	Type     *AliasPatternType `json:"type,omitempty"`
	Variable *string           `json:"variable,omitempty"`
}
