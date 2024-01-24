package volumequotarules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeQuotaRulesProperties struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	QuotaSizeInKiBs   *int64             `json:"quotaSizeInKiBs,omitempty"`
	QuotaTarget       *string            `json:"quotaTarget,omitempty"`
	QuotaType         *Type              `json:"quotaType,omitempty"`
}
