package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventInput struct {
	AccessControl            *LiveEventInputAccessControl `json:"accessControl,omitempty"`
	AccessToken              *string                      `json:"accessToken,omitempty"`
	Endpoints                *[]LiveEventEndpoint         `json:"endpoints,omitempty"`
	KeyFrameIntervalDuration *string                      `json:"keyFrameIntervalDuration,omitempty"`
	StreamingProtocol        LiveEventInputProtocol       `json:"streamingProtocol"`
}
