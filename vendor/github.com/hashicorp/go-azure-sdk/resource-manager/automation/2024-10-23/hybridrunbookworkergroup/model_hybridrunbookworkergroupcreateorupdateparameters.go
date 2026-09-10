package hybridrunbookworkergroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridRunbookWorkerGroupCreateOrUpdateParameters struct {
	Name       *string                                           `json:"name,omitempty"`
	Properties *HybridRunbookWorkerGroupCreateOrUpdateProperties `json:"properties,omitempty"`
}
