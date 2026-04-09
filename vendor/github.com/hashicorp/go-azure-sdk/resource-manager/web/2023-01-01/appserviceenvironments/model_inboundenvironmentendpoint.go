package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundEnvironmentEndpoint struct {
	Description *string   `json:"description,omitempty"`
	Endpoints   *[]string `json:"endpoints,omitempty"`
	Ports       *[]string `json:"ports,omitempty"`
}
