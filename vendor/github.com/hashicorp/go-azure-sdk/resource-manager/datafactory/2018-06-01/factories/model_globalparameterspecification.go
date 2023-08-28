package factories

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalParameterSpecification struct {
	Type  GlobalParameterType `json:"type"`
	Value interface{}         `json:"value"`
}
