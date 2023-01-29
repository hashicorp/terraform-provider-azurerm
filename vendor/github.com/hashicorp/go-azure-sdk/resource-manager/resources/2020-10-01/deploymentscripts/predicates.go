package deploymentscripts

type DeploymentScriptOperationPredicate struct {
}

func (p DeploymentScriptOperationPredicate) Matches(input DeploymentScript) bool {

	return true
}
