package findrestorabletimeranges

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupFindRestorableTimeRangesResponse struct {
	ObjectType           *string                `json:"objectType,omitempty"`
	RestorableTimeRanges *[]RestorableTimeRange `json:"restorableTimeRanges,omitempty"`
}
