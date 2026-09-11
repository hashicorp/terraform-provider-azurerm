//go:build tools

// Package tools pins the dev tool dependencies. The tag above means this file never compiles;
// the blank imports exist so go mod tidy treats the tools as direct dependencies and keeps them
// in their own require block, separate from the wall of transitive // indirect requirements.
package tools

import (
	_ "github.com/YakDriver/tfproviderdocs"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/katbyte/tctest"
	_ "github.com/katbyte/terrafmt"
	_ "github.com/rhysd/actionlint/cmd/actionlint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
