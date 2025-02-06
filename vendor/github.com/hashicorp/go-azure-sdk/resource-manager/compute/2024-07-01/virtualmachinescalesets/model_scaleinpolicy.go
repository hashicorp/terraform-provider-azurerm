package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleInPolicy struct {
	ForceDeletion *bool                                 `json:"forceDeletion,omitempty"`
	Rules         *[]VirtualMachineScaleSetScaleInRules `json:"rules,omitempty"`
}
