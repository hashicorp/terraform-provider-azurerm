package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAccount struct {
	ElevationLevel           *ElevationLevel           `json:"elevationLevel,omitempty"`
	LinuxUserConfiguration   *LinuxUserConfiguration   `json:"linuxUserConfiguration,omitempty"`
	Name                     string                    `json:"name"`
	Password                 string                    `json:"password"`
	WindowsUserConfiguration *WindowsUserConfiguration `json:"windowsUserConfiguration,omitempty"`
}
