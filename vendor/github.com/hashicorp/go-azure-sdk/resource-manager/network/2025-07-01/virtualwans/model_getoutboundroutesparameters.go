package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOutboundRoutesParameters struct {
	ConnectionType *string `json:"connectionType,omitempty"`
	ResourceUri    *string `json:"resourceUri,omitempty"`
}
