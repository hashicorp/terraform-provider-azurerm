package applications

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationPatchable struct {
	Id         *string                `json:"id,omitempty"`
	Identity   *Identity              `json:"identity,omitempty"`
	Kind       *string                `json:"kind,omitempty"`
	Location   *string                `json:"location,omitempty"`
	ManagedBy  *string                `json:"managedBy,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Plan       *PlanPatchable         `json:"plan,omitempty"`
	Properties *ApplicationProperties `json:"properties,omitempty"`
	Sku        *Sku                   `json:"sku,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
