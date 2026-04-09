package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteCloneability struct {
	BlockingCharacteristics *[]SiteCloneabilityCriterion `json:"blockingCharacteristics,omitempty"`
	BlockingFeatures        *[]SiteCloneabilityCriterion `json:"blockingFeatures,omitempty"`
	Result                  *CloneAbilityResult          `json:"result,omitempty"`
	UnsupportedFeatures     *[]SiteCloneabilityCriterion `json:"unsupportedFeatures,omitempty"`
}
