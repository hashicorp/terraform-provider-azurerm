package attestationproviders

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttestationProviderListResult struct {
	SystemData *systemdata.SystemData  `json:"systemData,omitempty"`
	Value      *[]AttestationProviders `json:"value,omitempty"`
}
