package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPAddress struct {
	Address    *string `json:"address,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`
}
