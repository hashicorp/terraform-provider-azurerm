package availabilitygrouplisteners

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgReplica struct {
	Commit                      *Commit            `json:"commit,omitempty"`
	Failover                    *Failover          `json:"failover,omitempty"`
	ReadableSecondary           *ReadableSecondary `json:"readableSecondary,omitempty"`
	Role                        *Role              `json:"role,omitempty"`
	SqlVirtualMachineInstanceId *string            `json:"sqlVirtualMachineInstanceId,omitempty"`
}
