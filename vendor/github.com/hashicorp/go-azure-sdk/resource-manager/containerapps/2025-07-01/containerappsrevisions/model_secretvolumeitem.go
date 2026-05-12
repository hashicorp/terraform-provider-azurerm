package containerappsrevisions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretVolumeItem struct {
	Path      *string `json:"path,omitempty"`
	SecretRef *string `json:"secretRef,omitempty"`
}
