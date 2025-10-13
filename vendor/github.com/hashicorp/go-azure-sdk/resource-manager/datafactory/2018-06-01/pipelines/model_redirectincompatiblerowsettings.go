package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedirectIncompatibleRowSettings struct {
	LinkedServiceName interface{}  `json:"linkedServiceName"`
	Path              *interface{} `json:"path,omitempty"`
}
