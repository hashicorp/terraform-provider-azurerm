package main

import (
	tfpasses "github.com/bflad/tfproviderlint/passes"
	tfxpasses "github.com/bflad/tfproviderlint/xpasses"
	azurermpasses "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tools/azurermproviderlint/passes"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	var analyzers []*analysis.Analyzer
	analyzers = append(analyzers, tfpasses.AllChecks...)
	analyzers = append(analyzers, tfxpasses.AllChecks...)
	analyzers = append(analyzers, azurermpasses.AllChecks...)
	multichecker.Main(analyzers...)
}
