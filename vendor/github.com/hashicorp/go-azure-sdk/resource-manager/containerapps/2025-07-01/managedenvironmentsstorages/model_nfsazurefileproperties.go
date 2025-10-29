package managedenvironmentsstorages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NfsAzureFileProperties struct {
	AccessMode *AccessMode `json:"accessMode,omitempty"`
	Server     *string     `json:"server,omitempty"`
	ShareName  *string     `json:"shareName,omitempty"`
}
