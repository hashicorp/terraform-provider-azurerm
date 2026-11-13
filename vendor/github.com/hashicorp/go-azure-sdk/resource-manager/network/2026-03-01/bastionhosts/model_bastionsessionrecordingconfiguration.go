package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionSessionRecordingConfiguration struct {
	BlobContainerUri string                   `json:"blobContainerUri"`
	Identity         SessionRecordingIdentity `json:"identity"`
}
