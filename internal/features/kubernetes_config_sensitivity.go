package features

import (
	"os"
	"strings"
)

func KubeConfigsAreSensitive() bool {
	return ThreePointOhBeta() || strings.EqualFold(os.Getenv("ARM_AKS_KUBE_CONFIGS_SENSITIVE"), "true")
}
