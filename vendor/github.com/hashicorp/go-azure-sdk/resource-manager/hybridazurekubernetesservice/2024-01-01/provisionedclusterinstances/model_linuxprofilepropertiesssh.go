package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxProfilePropertiesSsh struct {
	PublicKeys *[]LinuxProfilePropertiesSshPublicKeysInlined `json:"publicKeys,omitempty"`
}
