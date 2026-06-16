package containers

import (
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// suppressKubernetesVersionDiff suppresses the noisy no-op diff produced when an AKS
// automatic patch upgrade rewrites a configured minor version alias (e.g. "1.27") to the
// full patch version (e.g. "1.27.14"). When the configured value is a "major.minor" alias
// whose major and minor match the running patch version, the two are semantically equal
// (the alias resolves to the latest GA patch of that minor version, which is what is
// already running) so the diff is suppressed.
//
// It deliberately does not suppress when the configured value is a full patch version, so
// pinning an exact patch (e.g. "1.27.15" while running "1.27.14") still surfaces a real
// change.
func suppressKubernetesVersionDiff(_, old, new string, _ *pluginsdk.ResourceData) bool {
	// new = configured value, old = value currently in state.
	// Only a two-part "major.minor" configured value is treated as a minor version alias.
	configParts := strings.Split(new, ".")
	if len(configParts) != 2 {
		return false
	}

	stateParts := strings.Split(old, ".")
	if len(stateParts) < 2 {
		return false
	}

	return configParts[0] == stateParts[0] && configParts[1] == stateParts[1]
}
