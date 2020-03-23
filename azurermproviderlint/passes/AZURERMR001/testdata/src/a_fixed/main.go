package a

import (
	"errors"
	"fmt"
)

func f() {
	/* Passing case */

	_ = fmt.Errorf("failed to")
	_ = errors.New("failed to")

	/* Failing cases */

	_ = fmt.Errorf("failed to") // want `AZURERMR001: prefer other leading words instead of "error" as error message`
	_ = errors.New("failed to") // want `AZURERMR001: prefer other leading words instead of "error" as error message`

	/* Comment ignored cases */

	// lintignore:AZURERMR001
	_ = fmt.Errorf("Error to")
	_ = errors.New("error to") // lintignore:AZURERMR001
}
