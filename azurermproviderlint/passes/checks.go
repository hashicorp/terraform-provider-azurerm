package passes

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurermproviderlint/passes/AZURERMR001"
	"golang.org/x/tools/go/analysis"
)

var AllChecks = []*analysis.Analyzer{
	AZURERMR001.Analyzer,
}
