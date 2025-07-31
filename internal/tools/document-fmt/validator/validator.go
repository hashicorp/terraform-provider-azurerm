package validator

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/rule"
)

// Validator.Run provides a separation between `data` and `rule` packages to avoid cyclic imports...
// TODO: review tool structure?

type Validator struct{}

func (v Validator) Run(rules []string, data *data.TerraformNodeData, fix bool) {
	for _, id := range rules {
		if r := rule.GetRule(id); r != nil {
			errs := r.Run(data, fix)
			if len(errs) > 0 {
				data.Errors = append(data.Errors, errs...)
			}
		}
	}
}
