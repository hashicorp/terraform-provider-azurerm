package connectedregistries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyncUpdateProperties struct {
	MessageTtl *string `json:"messageTtl,omitempty"`
	Schedule   *string `json:"schedule,omitempty"`
	SyncWindow *string `json:"syncWindow,omitempty"`
}
