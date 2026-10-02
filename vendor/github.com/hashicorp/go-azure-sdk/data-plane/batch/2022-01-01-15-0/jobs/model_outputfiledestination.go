package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutputFileDestination struct {
	Container *OutputFileBlobContainerDestination `json:"container,omitempty"`
}
