package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TemplateHashResult struct {
	MinifiedTemplate *string `json:"minifiedTemplate,omitempty"`
	TemplateHash     *string `json:"templateHash,omitempty"`
}
