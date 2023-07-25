package devices

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataBoxEdgeDevice struct {
	Etag       *string                      `json:"etag,omitempty"`
	Id         *string                      `json:"id,omitempty"`
	Identity   *ResourceIdentity            `json:"identity,omitempty"`
	Kind       *DataBoxEdgeDeviceKind       `json:"kind,omitempty"`
	Location   string                       `json:"location"`
	Name       *string                      `json:"name,omitempty"`
	Properties *DataBoxEdgeDeviceProperties `json:"properties,omitempty"`
	Sku        *Sku                         `json:"sku,omitempty"`
	SystemData *systemdata.SystemData       `json:"systemData,omitempty"`
	Tags       *map[string]string           `json:"tags,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}
