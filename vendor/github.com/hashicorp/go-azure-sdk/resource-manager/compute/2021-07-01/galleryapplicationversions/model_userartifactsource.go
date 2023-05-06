package galleryapplicationversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserArtifactSource struct {
	DefaultConfigurationLink *string `json:"defaultConfigurationLink,omitempty"`
	MediaLink                string  `json:"mediaLink"`
}
