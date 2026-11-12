package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PeerDbDetails struct {
	PeerDbId       *string `json:"peerDbId,omitempty"`
	PeerDbLocation *string `json:"peerDbLocation,omitempty"`
	PeerDbOcid     *string `json:"peerDbOcid,omitempty"`
}
