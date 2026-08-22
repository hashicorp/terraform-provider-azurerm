package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HardwareProfile struct {
	NumberOfCPUSockets         *int64       `json:"numberOfCpuSockets,omitempty"`
	Processors                 *[]Processor `json:"processors,omitempty"`
	TotalPhysicalMemoryInBytes *int64       `json:"totalPhysicalMemoryInBytes,omitempty"`
}
