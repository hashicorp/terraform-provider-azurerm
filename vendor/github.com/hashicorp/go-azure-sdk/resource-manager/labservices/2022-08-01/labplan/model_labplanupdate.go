package labplan

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabPlanUpdate struct {
	Identity   *identity.SystemAssigned `json:"identity,omitempty"`
	Properties *LabPlanUpdateProperties `json:"properties,omitempty"`
	Tags       *[]string                `json:"tags,omitempty"`
}
