package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPrivateEndpointConnection struct {
	Id         *string                              `json:"id,omitempty"`
	Properties *PrivateEndpointConnectionProperties `json:"properties,omitempty"`
}
