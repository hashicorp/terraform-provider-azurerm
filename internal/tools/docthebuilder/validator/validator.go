package validator

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/rule"
)

// Validator.Run provides a separation between `data` and `rule` packages to avoid cyclic imports...
// TODO: review tool structure?

type Validator struct{}

func (v Validator) Run(rules []string, data *data.ResourceData, fix bool) {
	for _, name := range rules {
		if r := rule.GetRule(name); r != nil {
			errs := r.Run(data, fix)
			if len(errs) > 0 {
				data.Errors = append(data.Errors, errs...)
			}
		}
	}
}
