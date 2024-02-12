package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImmutabilitySettings struct {
	State *ImmutabilityState `json:"state,omitempty"`
}
