package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPath struct {
	EncryptionScheme  EncryptionScheme                 `json:"encryptionScheme"`
	Paths             *[]string                        `json:"paths,omitempty"`
	StreamingProtocol StreamingPolicyStreamingProtocol `json:"streamingProtocol"`
}
