package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DatabaseCollation() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`(^[A-Z]+)([A-Za-z0-9]+_)+((BIN|BIN2|CI_AI|CI_AI_KS|CI_AI_KS_WS|CI_AI_WS|CI_AS|CI_AS_KS|CI_AS_KS_WS|CS_AI|CS_AI_KS|CS_AI_KS_WS|CS_AI_WS|CS_AS|CS_AS_KS|CS_AS_KS_WS|CS_AS_WS)+)((_[A-Za-z0-9]+)+$)*`),

		`This is not a valid collation.`,
	)
}
