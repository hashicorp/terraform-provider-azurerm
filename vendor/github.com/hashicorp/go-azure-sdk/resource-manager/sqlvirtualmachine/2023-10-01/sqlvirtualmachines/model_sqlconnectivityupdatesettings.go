package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlConnectivityUpdateSettings struct {
	ConnectivityType      *ConnectivityType `json:"connectivityType,omitempty"`
	Port                  *int64            `json:"port,omitempty"`
	SqlAuthUpdatePassword *string           `json:"sqlAuthUpdatePassword,omitempty"`
	SqlAuthUpdateUserName *string           `json:"sqlAuthUpdateUserName,omitempty"`
}
