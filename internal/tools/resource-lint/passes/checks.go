package passes

import (
	"golang.org/x/tools/go/analysis"
)

// AllChecks contains all Analyzers that report issues
var AllChecks = []*analysis.Analyzer{
	AZBP001Analyzer,
	AZBP002Analyzer,
	AZBP003Analyzer,
	AZBP004Analyzer,
	AZSD001Analyzer,
	AZSD002Analyzer,
	AZRN001Analyzer,
	AZRE001Analyzer,
	AZNR001Analyzer,
}
