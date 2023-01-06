package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventTranscription struct {
	InputTrackSelection      *[]LiveEventInputTrackSelection    `json:"inputTrackSelection,omitempty"`
	Language                 *string                            `json:"language,omitempty"`
	OutputTranscriptionTrack *LiveEventOutputTranscriptionTrack `json:"outputTranscriptionTrack,omitempty"`
}
