package providerinstances

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderInstance struct {
	Id         *string                     `json:"id,omitempty"`
	Identity   *identity.UserAssignedMap   `json:"identity,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	Properties *ProviderInstanceProperties `json:"properties,omitempty"`
	SystemData *systemdata.SystemData      `json:"systemData,omitempty"`
	Type       *string                     `json:"type,omitempty"`
}
