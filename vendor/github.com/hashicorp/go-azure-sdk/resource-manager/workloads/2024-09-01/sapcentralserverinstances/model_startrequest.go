package sapcentralserverinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartRequest struct {
	StartVM *bool `json:"startVm,omitempty"`
}
