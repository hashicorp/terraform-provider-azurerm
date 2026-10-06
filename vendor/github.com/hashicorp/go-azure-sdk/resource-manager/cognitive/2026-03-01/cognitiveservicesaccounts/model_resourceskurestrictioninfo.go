package cognitiveservicesaccounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSkuRestrictionInfo struct {
	Locations *[]string     `json:"locations,omitempty"`
	Zones     *zones.Schema `json:"zones,omitempty"`
}
