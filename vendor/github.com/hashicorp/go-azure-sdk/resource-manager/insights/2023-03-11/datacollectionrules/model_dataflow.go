package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataFlow struct {
	BuiltInTransform *string                 `json:"builtInTransform,omitempty"`
	CaptureOverflow  *bool                   `json:"captureOverflow,omitempty"`
	Destinations     *[]string               `json:"destinations,omitempty"`
	OutputStream     *string                 `json:"outputStream,omitempty"`
	Streams          *[]KnownDataFlowStreams `json:"streams,omitempty"`
	TransformKql     *string                 `json:"transformKql,omitempty"`
}
