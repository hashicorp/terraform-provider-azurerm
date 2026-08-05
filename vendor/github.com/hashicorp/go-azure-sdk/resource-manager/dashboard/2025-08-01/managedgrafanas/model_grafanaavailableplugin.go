package managedgrafanas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GrafanaAvailablePlugin struct {
	Author   *string `json:"author,omitempty"`
	Name     *string `json:"name,omitempty"`
	PluginId *string `json:"pluginId,omitempty"`
	Type     *string `json:"type,omitempty"`
}
