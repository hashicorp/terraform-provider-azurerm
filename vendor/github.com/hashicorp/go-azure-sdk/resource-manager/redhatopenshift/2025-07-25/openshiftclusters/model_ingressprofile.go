package openshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngressProfile struct {
	IP         *string     `json:"ip,omitempty"`
	Name       *string     `json:"name,omitempty"`
	Visibility *Visibility `json:"visibility,omitempty"`
}
