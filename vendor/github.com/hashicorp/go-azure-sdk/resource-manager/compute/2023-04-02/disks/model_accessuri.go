package disks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessUri struct {
	AccessSAS             *string `json:"accessSAS,omitempty"`
	SecurityDataAccessSAS *string `json:"securityDataAccessSAS,omitempty"`
}
