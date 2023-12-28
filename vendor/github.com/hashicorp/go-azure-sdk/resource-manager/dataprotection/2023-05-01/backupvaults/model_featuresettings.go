package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FeatureSettings struct {
	CrossRegionRestoreSettings       *CrossRegionRestoreSettings       `json:"crossRegionRestoreSettings,omitempty"`
	CrossSubscriptionRestoreSettings *CrossSubscriptionRestoreSettings `json:"crossSubscriptionRestoreSettings,omitempty"`
}
