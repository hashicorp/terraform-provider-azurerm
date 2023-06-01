package galleries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommunityGalleryInfo struct {
	CommunityGalleryEnabled *bool     `json:"communityGalleryEnabled,omitempty"`
	Eula                    *string   `json:"eula,omitempty"`
	PublicNamePrefix        *string   `json:"publicNamePrefix,omitempty"`
	PublicNames             *[]string `json:"publicNames,omitempty"`
	PublisherContact        *string   `json:"publisherContact,omitempty"`
	PublisherUri            *string   `json:"publisherUri,omitempty"`
}
