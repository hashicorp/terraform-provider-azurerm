package differror

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
)

type DiffError struct {
	t string
	c string
	e string
}

func New(text, current, expected string) error {
	return &DiffError{
		t: text,
		c: current,
		e: expected,
	}
}

func (e DiffError) Error() string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("%s\n", e.t))
	b.WriteString(util.Red(fmt.Sprintf("- %s\n", e.c)))
	b.WriteString(util.Green(fmt.Sprintf("+ %s", e.e)))

	return b.String()
}
