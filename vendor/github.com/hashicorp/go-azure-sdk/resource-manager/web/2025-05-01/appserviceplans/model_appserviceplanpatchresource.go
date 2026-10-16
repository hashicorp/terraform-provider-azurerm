package appserviceplans

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServicePlanPatchResource struct {
	Id         *string                                `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap     `json:"identity,omitempty"`
	Kind       *string                                `json:"kind,omitempty"`
	Name       *string                                `json:"name,omitempty"`
	Properties *AppServicePlanPatchResourceProperties `json:"properties,omitempty"`
	Type       *string                                `json:"type,omitempty"`
}
