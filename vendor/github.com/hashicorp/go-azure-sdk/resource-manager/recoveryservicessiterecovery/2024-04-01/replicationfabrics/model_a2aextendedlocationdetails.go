package replicationfabrics

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AExtendedLocationDetails struct {
	PrimaryExtendedLocation  *edgezones.Model `json:"primaryExtendedLocation,omitempty"`
	RecoveryExtendedLocation *edgezones.Model `json:"recoveryExtendedLocation,omitempty"`
}
