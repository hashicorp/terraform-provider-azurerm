package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityContextCapabilitiesDefinition struct {
	Add  *[]string `json:"add,omitempty"`
	Drop *[]string `json:"drop,omitempty"`
}
