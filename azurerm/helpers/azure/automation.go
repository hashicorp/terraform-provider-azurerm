package azure

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ValidateAutomationAccountName validates Automation Account names
func ValidateAutomationAccountName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z]{4,48}[0-9a-zA-Z]$`),
		`The account name must start with a letter or number.  The account name can contain letters, numbers, and dashes. The final character must be a letter or a number. The account name length must be from 6 to 50 characters.`,
	)
}

// ValidateAutomationRunbookName validates Automation Account Runbook names
func ValidateAutomationRunbookName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[0-9a-zA-Z][-_0-9a-zA-Z]{0,62}$`),
		`The name can contain only letters, numbers, underscores and dashes. The name must begin with a letter. The name must be less than 64 characters.`,
	)
}

// ValidateAutomationScheduleName validates Automation Account Schedule names
func ValidateAutomationScheduleName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[^<>*%&:\\?.+/]{0,127}[^<>*%&:\\?.+/\s]$`),
		`The name length must be from 1 to 128 characters. The name cannot contain special characters < > * % & : \ ? . + / and cannot end with a whitespace character.`,
	)
}
