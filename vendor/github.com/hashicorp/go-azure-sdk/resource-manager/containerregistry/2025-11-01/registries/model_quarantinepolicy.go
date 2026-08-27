package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuarantinePolicy struct {
	Status *PolicyStatus `json:"status,omitempty"`
}
