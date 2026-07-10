package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// CallTerraform calls the locally installed terraform binary with the specified options
// returns the raw []byte output for the caller to process.
func CallTerraform(opts ...string) ([]byte, error) {
	cmd := exec.Command("terraform", opts...)
	out, err := cmd.Output()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("missing terraform binary, please ensure Terraform is installed and discoverable in your PATH")
		}
		return nil, err
	}

	return out, nil
}

// GoFmt calls `gofmt -w` over the specified file (including path)
func GoFmt(file string) error {
	cmd := exec.Command("gofmt", "-w", fmt.Sprintf("./%s", file))
	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}

// GoImports runs `goimports -w` over the specified file to fix and prune the
// import block, falling back to `gofmt -w` when goimports is not installed.
func GoImports(file string) error {
	cmd := exec.Command("goimports", "-w", file)
	if out, err := cmd.CombinedOutput(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return GoFmt(file)
		}
		return fmt.Errorf("running goimports on %s: %v: %s", file, err, string(out))
	}
	return nil
}

// Terrafmt calls (if installed) katbyte/terrafmt to format Terraform
// configurations in the specified file
func Terrafmt(path string) error {
	cmd := exec.Command("terrafmt", "fmt", "-f", path)
	if _, err := cmd.Output(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("missing terrafmt, please ensure katbyte/terrafmt is installed (`go install github.com/katbyte/terrafmt@latest`) and discoverable in your PATH")
		}
		return err
	}

	return nil
}
