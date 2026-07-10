package ir

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

// acronymPluralRe matches an all-caps acronym immediately followed by a trailing
// lowercase "s" at a word boundary (e.g. IPs, URLs, IDs, VMs).
var acronymPluralRe = regexp.MustCompile(`([A-Z][A-Z]+)s\b`)

// normalizeAcronymPlurals rewrites pluralised acronyms (IPs -> Ips, URLs ->
// Urls) so snake casing yields ips/urls rather than i_ps/u_rls, matching the
// provider's hand-written naming.
func normalizeAcronymPlurals(s string) string {
	return acronymPluralRe.ReplaceAllStringFunc(s, func(m string) string {
		body := m[:len(m)-1] // strip trailing 's'
		return body[:1] + strings.ToLower(body[1:]) + "s"
	})
}

// snake converts a PascalCase/camelCase identifier to snake_case. strcase
// handles acronyms correctly (e.g. VMSize -> vm_size, PreconfiguredNSG ->
// preconfigured_nsg, IP -> ip), matching the provider's hand-written naming.
func snake(s string) string {
	return strcase.ToSnake(normalizeAcronymPlurals(s))
}

// camel converts an identifier to PascalCase.
func camel(s string) string {
	return strcase.ToCamel(s)
}

// idToID normalises a trailing "Id" to "ID" so generated Parse/New function
// names match the go-azure-sdk convention (e.g. OpenShiftClusterId ->
// OpenShiftClusterID, giving ParseOpenShiftClusterID / NewOpenShiftClusterID).
func idToID(id string) string {
	if strings.HasSuffix(id, "Id") {
		return strings.TrimSuffix(id, "Id") + "ID"
	}
	return id
}

// singularize returns a best-effort singular form of a resource name, used for
// deriving Go identifiers from plural Pandora resource keys (e.g.
// OpenShiftClusters -> OpenShiftCluster).
func singularize(s string) string {
	switch {
	case strings.HasSuffix(s, "ies"):
		return strings.TrimSuffix(s, "ies") + "y"
	case strings.HasSuffix(s, "sses"):
		return strings.TrimSuffix(s, "es")
	case strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss"):
		return strings.TrimSuffix(s, "s")
	default:
		return s
	}
}
