package encodings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PresetConfigurations struct {
	Complexity                *Complexity       `json:"complexity,omitempty"`
	InterleaveOutput          *InterleaveOutput `json:"interleaveOutput,omitempty"`
	KeyFrameIntervalInSeconds *float64          `json:"keyFrameIntervalInSeconds,omitempty"`
	MaxBitrateBps             *int64            `json:"maxBitrateBps,omitempty"`
	MaxHeight                 *int64            `json:"maxHeight,omitempty"`
	MaxLayers                 *int64            `json:"maxLayers,omitempty"`
	MinBitrateBps             *int64            `json:"minBitrateBps,omitempty"`
	MinHeight                 *int64            `json:"minHeight,omitempty"`
}
