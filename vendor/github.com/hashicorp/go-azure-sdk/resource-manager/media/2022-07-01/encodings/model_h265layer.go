package encodings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type H265Layer struct {
	AdaptiveBFrame  *bool             `json:"adaptiveBFrame,omitempty"`
	BFrames         *int64            `json:"bFrames,omitempty"`
	Bitrate         int64             `json:"bitrate"`
	BufferWindow    *string           `json:"bufferWindow,omitempty"`
	Crf             *float64          `json:"crf,omitempty"`
	FrameRate       *string           `json:"frameRate,omitempty"`
	Height          *string           `json:"height,omitempty"`
	Label           *string           `json:"label,omitempty"`
	Level           *string           `json:"level,omitempty"`
	MaxBitrate      *int64            `json:"maxBitrate,omitempty"`
	Profile         *H265VideoProfile `json:"profile,omitempty"`
	ReferenceFrames *int64            `json:"referenceFrames,omitempty"`
	Slices          *int64            `json:"slices,omitempty"`
	Width           *string           `json:"width,omitempty"`
}
