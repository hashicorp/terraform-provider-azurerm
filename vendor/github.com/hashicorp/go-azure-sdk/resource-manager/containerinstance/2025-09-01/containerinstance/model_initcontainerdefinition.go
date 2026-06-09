package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InitContainerDefinition struct {
	Name       string                            `json:"name"`
	Properties InitContainerPropertiesDefinition `json:"properties"`
}
