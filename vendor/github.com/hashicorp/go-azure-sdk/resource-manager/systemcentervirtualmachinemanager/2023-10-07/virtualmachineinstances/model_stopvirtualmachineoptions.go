package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StopVirtualMachineOptions struct {
	SkipShutdown *SkipShutdown `json:"skipShutdown,omitempty"`
}
