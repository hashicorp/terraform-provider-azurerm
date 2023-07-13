package deploymentscripts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentScriptOperationPredicate struct {
}

func (p DeploymentScriptOperationPredicate) Matches(input DeploymentScript) bool {

	return true
}
