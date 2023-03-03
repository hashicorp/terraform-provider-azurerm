package alertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FusionSourceSettings struct {
	Enabled        bool                          `json:"enabled"`
	SourceName     string                        `json:"sourceName"`
	SourceSubTypes *[]FusionSourceSubTypeSetting `json:"sourceSubTypes,omitempty"`
}
