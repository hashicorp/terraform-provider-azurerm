package assetsandassetfilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HlsSettings struct {
	Characteristics *string `json:"characteristics,omitempty"`
	Default         *bool   `json:"default,omitempty"`
	Forced          *bool   `json:"forced,omitempty"`
}
