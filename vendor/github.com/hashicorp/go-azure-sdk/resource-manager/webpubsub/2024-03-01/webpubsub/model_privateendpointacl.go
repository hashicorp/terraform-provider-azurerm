package webpubsub

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointACL struct {
	Allow *[]WebPubSubRequestType `json:"allow,omitempty"`
	Deny  *[]WebPubSubRequestType `json:"deny,omitempty"`
	Name  string                  `json:"name"`
}
