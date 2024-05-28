package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineCaptureResult struct {
	ContentVersion *string        `json:"contentVersion,omitempty"`
	Id             *string        `json:"id,omitempty"`
	Parameters     *interface{}   `json:"parameters,omitempty"`
	Resources      *[]interface{} `json:"resources,omitempty"`
	Schema         *string        `json:"$schema,omitempty"`
}
