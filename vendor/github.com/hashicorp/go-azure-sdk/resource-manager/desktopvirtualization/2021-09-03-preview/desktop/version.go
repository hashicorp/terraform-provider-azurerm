package desktop

import "fmt"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

const defaultApiVersion = "2021-09-03-preview"

func userAgent() string {
	return fmt.Sprintf("hashicorp/go-azure-sdk/desktop/%s", defaultApiVersion)
}
