package validate

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func FirewallPolicyName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_][\w-.]*[\w]$`),
		"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.")
}

func FirewallPolicyRuleCollectionGroupName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_][\w-.]*[\w]$`),
		"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.")
}

func FirewallPolicyRuleName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^\W_][\w-.]*[\w]$`),
		"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.")
}

func FirewallPolicyID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.FirewallPolicyID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func FirewallPolicyRuleCollectionGroupID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.FirewallPolicyRuleCollectionGroupID(v); err != nil {
		errors = append(errors, fmt.Errorf("parsing %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func FirewallPolicyRulePort(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	assertWithinRnage := func(n int) error {
		if n < 1 || n > 65535 {
			return fmt.Errorf("port %d is out of range (1-65535)", n)
		}
		return nil
	}

	// Allowed format including: `num` or `num1-num2` (num1 < num2).
	groups := regexp.MustCompile(`^(\d+)((-)(\d+))?$`).FindStringSubmatch(v)
	if len(groups) != 5 {
		errors = append(errors, fmt.Errorf("invalid format of %q", k))
		return
	}
	if groups[2] == "" {
		p1, _ := strconv.Atoi(groups[1])
		if err := assertWithinRnage(p1); err != nil {
			errors = append(errors, err)
			return
		}
	} else {
		p1, _ := strconv.Atoi(groups[1])
		p2, _ := strconv.Atoi(groups[4])
		if p1 >= p2 {
			errors = append(errors, fmt.Errorf("beginning port (%d) should be less than endping port (%d)", p1, p2))
			return
		}
		if err := assertWithinRnage(p1); err != nil {
			errors = append(errors, err)
			return
		}
		if err := assertWithinRnage(p2); err != nil {
			errors = append(errors, err)
			return
		}
	}

	return nil, nil
}
