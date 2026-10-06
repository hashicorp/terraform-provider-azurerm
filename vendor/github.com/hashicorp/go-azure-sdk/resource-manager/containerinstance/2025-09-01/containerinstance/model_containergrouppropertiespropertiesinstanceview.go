package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupPropertiesPropertiesInstanceView struct {
	Events *[]Event `json:"events,omitempty"`
	State  *string  `json:"state,omitempty"`
}
