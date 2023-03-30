package contentkeypolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyFairPlayOfflineRentalConfiguration struct {
	PlaybackDurationSeconds int64 `json:"playbackDurationSeconds"`
	StorageDurationSeconds  int64 `json:"storageDurationSeconds"`
}
