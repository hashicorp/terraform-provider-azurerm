package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateProfile struct {
	RollingUpdateProfile *UpdateProfileRollingUpdateProfile `json:"rollingUpdateProfile,omitempty"`
	UpdateMode           *NGroupUpdateMode                  `json:"updateMode,omitempty"`
}
