package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectorOperationPredicate struct {
}

func (p DataConnectorOperationPredicate) Matches(input DataConnector) bool {

	return true
}
