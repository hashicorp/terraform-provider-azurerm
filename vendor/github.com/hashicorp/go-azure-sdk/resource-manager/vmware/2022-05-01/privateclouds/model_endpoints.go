package privateclouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Endpoints struct {
	HcxCloudManager *string `json:"hcxCloudManager,omitempty"`
	NsxtManager     *string `json:"nsxtManager,omitempty"`
	Vcsa            *string `json:"vcsa,omitempty"`
}
