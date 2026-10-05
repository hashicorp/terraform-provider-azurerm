package standbyvirtualmachinepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StandbyVirtualMachinePoolElasticityProfile struct {
	MaxReadyCapacity int64  `json:"maxReadyCapacity"`
	MinReadyCapacity *int64 `json:"minReadyCapacity,omitempty"`
}
