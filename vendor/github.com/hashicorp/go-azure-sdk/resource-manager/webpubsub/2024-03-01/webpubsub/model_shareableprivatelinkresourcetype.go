package webpubsub

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ShareablePrivateLinkResourceType struct {
	Name       *string                                 `json:"name,omitempty"`
	Properties *ShareablePrivateLinkResourceProperties `json:"properties,omitempty"`
}
