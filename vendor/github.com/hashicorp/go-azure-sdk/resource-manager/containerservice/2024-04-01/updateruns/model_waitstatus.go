package updateruns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WaitStatus struct {
	Status                *UpdateStatus `json:"status,omitempty"`
	WaitDurationInSeconds *int64        `json:"waitDurationInSeconds,omitempty"`
}
