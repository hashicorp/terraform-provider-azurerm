package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateAccessDefinition struct {
	Approver  *string            `json:"approver,omitempty"`
	Metadata  JitRequestMetadata `json:"metadata"`
	Status    Status             `json:"status"`
	SubStatus Substatus          `json:"subStatus"`
}
