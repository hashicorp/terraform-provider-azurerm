package containerappsrevisions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerAppProbeTcpSocket struct {
	Host *string `json:"host,omitempty"`
	Port int64   `json:"port"`
}
