package videoanalyzer

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoProperties struct {
	Description *string         `json:"description,omitempty"`
	Flags       *VideoFlags     `json:"flags,omitempty"`
	MediaInfo   *VideoMediaInfo `json:"mediaInfo,omitempty"`
	Streaming   *VideoStreaming `json:"streaming,omitempty"`
	Title       *string         `json:"title,omitempty"`
	Type        *VideoType      `json:"type,omitempty"`
}
