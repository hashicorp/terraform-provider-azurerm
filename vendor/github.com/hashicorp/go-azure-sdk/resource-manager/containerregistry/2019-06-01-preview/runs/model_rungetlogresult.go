package runs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunGetLogResult struct {
	LogArtifactLink *string `json:"logArtifactLink,omitempty"`
	LogLink         *string `json:"logLink,omitempty"`
}
