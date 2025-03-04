package replicationfabrics

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AFabricSpecificLocationDetails struct {
	InitialPrimaryExtendedLocation  *edgezones.Model `json:"initialPrimaryExtendedLocation,omitempty"`
	InitialPrimaryFabricLocation    *string          `json:"initialPrimaryFabricLocation,omitempty"`
	InitialPrimaryZone              *string          `json:"initialPrimaryZone,omitempty"`
	InitialRecoveryExtendedLocation *edgezones.Model `json:"initialRecoveryExtendedLocation,omitempty"`
	InitialRecoveryFabricLocation   *string          `json:"initialRecoveryFabricLocation,omitempty"`
	InitialRecoveryZone             *string          `json:"initialRecoveryZone,omitempty"`
	PrimaryExtendedLocation         *edgezones.Model `json:"primaryExtendedLocation,omitempty"`
	PrimaryFabricLocation           *string          `json:"primaryFabricLocation,omitempty"`
	PrimaryZone                     *string          `json:"primaryZone,omitempty"`
	RecoveryExtendedLocation        *edgezones.Model `json:"recoveryExtendedLocation,omitempty"`
	RecoveryFabricLocation          *string          `json:"recoveryFabricLocation,omitempty"`
	RecoveryZone                    *string          `json:"recoveryZone,omitempty"`
}
