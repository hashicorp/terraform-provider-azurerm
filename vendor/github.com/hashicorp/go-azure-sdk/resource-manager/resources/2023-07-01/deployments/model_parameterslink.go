package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParametersLink struct {
	ContentVersion *string `json:"contentVersion,omitempty"`
	Uri            string  `json:"uri"`
}
