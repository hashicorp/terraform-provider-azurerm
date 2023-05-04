package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Source struct {
	Addr       *string `json:"addr,omitempty"`
	InstanceID *string `json:"instanceID,omitempty"`
}
