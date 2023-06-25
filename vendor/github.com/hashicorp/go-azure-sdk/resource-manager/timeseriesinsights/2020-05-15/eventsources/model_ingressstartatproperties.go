package eventsources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngressStartAtProperties struct {
	Time *string             `json:"time,omitempty"`
	Type *IngressStartAtType `json:"type,omitempty"`
}
