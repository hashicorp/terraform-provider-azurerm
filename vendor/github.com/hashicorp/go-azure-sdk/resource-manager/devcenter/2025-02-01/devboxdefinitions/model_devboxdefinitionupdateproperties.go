package devboxdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevBoxDefinitionUpdateProperties struct {
	HibernateSupport *HibernateSupport `json:"hibernateSupport,omitempty"`
	ImageReference   *ImageReference   `json:"imageReference,omitempty"`
	OsStorageType    *string           `json:"osStorageType,omitempty"`
	Sku              *Sku              `json:"sku,omitempty"`
}
