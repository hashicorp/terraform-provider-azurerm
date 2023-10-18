package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtendedLocation struct {
	Name *string                `json:"name,omitempty"`
	Type *ExtendedLocationTypes `json:"type,omitempty"`
}
