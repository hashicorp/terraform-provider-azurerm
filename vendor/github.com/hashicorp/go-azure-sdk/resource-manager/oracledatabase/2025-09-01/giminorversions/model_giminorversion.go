package giminorversions

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GiMinorVersion struct {
	Id         *string                   `json:"id,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *GiMinorVersionProperties `json:"properties,omitempty"`
	SystemData *systemdata.SystemData    `json:"systemData,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
