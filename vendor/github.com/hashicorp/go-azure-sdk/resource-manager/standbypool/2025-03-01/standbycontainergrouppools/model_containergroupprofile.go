package standbycontainergrouppools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupProfile struct {
	Id       string `json:"id"`
	Revision *int64 `json:"revision,omitempty"`
}
