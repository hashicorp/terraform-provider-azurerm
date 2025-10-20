package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenerateResourceLimits struct {
	Cpu *OperationalMode `json:"cpu,omitempty"`
}
