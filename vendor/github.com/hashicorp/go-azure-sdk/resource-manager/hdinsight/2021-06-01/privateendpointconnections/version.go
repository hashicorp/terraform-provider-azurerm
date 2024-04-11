package privateendpointconnections

import "fmt"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2021-06-01"

func userAgent() string {
	return fmt.Sprintf("hashicorp/go-azure-sdk/privateendpointconnections/%s", defaultApiVersion)
}
