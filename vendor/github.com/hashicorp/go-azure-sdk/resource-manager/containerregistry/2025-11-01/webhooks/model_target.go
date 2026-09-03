package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Target struct {
	Digest     *string `json:"digest,omitempty"`
	Length     *int64  `json:"length,omitempty"`
	MediaType  *string `json:"mediaType,omitempty"`
	Name       *string `json:"name,omitempty"`
	Repository *string `json:"repository,omitempty"`
	Size       *int64  `json:"size,omitempty"`
	Tag        *string `json:"tag,omitempty"`
	Url        *string `json:"url,omitempty"`
	Version    *string `json:"version,omitempty"`
}
