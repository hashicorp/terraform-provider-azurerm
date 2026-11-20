// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sdk

import "encoding/json"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceValidationRequestResource struct {
	ApiVersion       string             `json:"apiVersion"`
	Id               *string            `json:"id,omitempty"`
	Name             string             `json:"name"`
	Type             string             `json:"type"`
	Location         *string            `json:"location,omitempty"`
	Kind             *string            `json:"kind,omitempty"`
	ManagedBy        *string            `json:"managedBy,omitempty"`
	Tags             *map[string]string `json:"tags,omitempty"`
	Zones            *[]string          `json:"zones,omitempty"`
	ExtendedLocation json.RawMessage    `json:"extendedLocation,omitempty"`
	Identity         json.RawMessage    `json:"identity,omitempty"`
	Plan             json.RawMessage    `json:"plan,omitempty"`
	Properties       json.RawMessage    `json:"properties,omitempty"`
	Sku              json.RawMessage    `json:"sku,omitempty"`
}
