package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigServerSettingsErrorRecord struct {
	Messages *[]string `json:"messages,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Uri      *string   `json:"uri,omitempty"`
}
