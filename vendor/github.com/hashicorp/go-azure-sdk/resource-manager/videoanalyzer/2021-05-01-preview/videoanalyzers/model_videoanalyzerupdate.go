package videoanalyzers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoAnalyzerUpdate struct {
	Identity   *VideoAnalyzerIdentity         `json:"identity,omitempty"`
	Properties *VideoAnalyzerPropertiesUpdate `json:"properties,omitempty"`
	Tags       *map[string]string             `json:"tags,omitempty"`
}
