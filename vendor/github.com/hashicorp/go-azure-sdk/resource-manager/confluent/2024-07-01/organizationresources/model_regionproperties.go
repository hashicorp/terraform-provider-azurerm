package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionProperties struct {
	Metadata *SCMetadataEntity `json:"metadata,omitempty"`
	Spec     *RegionSpecEntity `json:"spec,omitempty"`
}
