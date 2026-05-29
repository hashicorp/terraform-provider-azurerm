package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineCaptureParameters struct {
	DestinationContainerName string `json:"destinationContainerName"`
	OverwriteVhds            bool   `json:"overwriteVhds"`
	VhdPrefix                string `json:"vhdPrefix"`
}
