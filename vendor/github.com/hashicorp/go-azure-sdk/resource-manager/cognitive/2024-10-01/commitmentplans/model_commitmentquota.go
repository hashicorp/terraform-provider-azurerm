package commitmentplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommitmentQuota struct {
	Quantity *int64  `json:"quantity,omitempty"`
	Unit     *string `json:"unit,omitempty"`
}
