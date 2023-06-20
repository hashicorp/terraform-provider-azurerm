package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventEncoding struct {
	EncodingType     *LiveEventEncodingType `json:"encodingType,omitempty"`
	KeyFrameInterval *string                `json:"keyFrameInterval,omitempty"`
	PresetName       *string                `json:"presetName,omitempty"`
	StretchMode      *StretchMode           `json:"stretchMode,omitempty"`
}
