// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package reporting

import (
	"fmt"
	"go/token"
	"sync"

	"golang.org/x/tools/go/analysis"
)

const (
	MatchModeExactAdded  = "exact-added"
	MatchModeSameHunk    = "same-hunk"
	MatchModeNewFile     = "new-file"
	MatchModeFileChanged = "file-changed"
)

type DiagnosticMeta struct {
	PkgPath       string
	Rule          string
	Message       string
	ReportFile    string
	ReportLine    int
	ReportColumn  int
	EvidenceFile  string
	EvidenceLines []int
	MatchMode     string
}

type ReportOptions struct {
	Rule          string
	ReportPos     token.Pos
	Message       string
	EvidenceFile  string
	EvidenceLines []int
	MatchMode     string
}

var (
	registryMu sync.Mutex
	registry   = make(map[string]DiagnosticMeta)
)

func Reset() {
	registryMu.Lock()
	defer registryMu.Unlock()

	registry = make(map[string]DiagnosticMeta)
}

func Record(meta DiagnosticMeta) {
	registryMu.Lock()
	defer registryMu.Unlock()

	meta = normalizeMeta(meta)
	registry[makeKey(meta.PkgPath, meta.ReportFile, meta.ReportLine, meta.ReportColumn, meta.Message)] = cloneMeta(meta)
}

func Lookup(pkgPath, file string, line, column int, message string) (DiagnosticMeta, bool) {
	registryMu.Lock()
	defer registryMu.Unlock()

	meta, ok := registry[makeKey(pkgPath, file, line, column, message)]
	if !ok {
		return DiagnosticMeta{}, false
	}

	return cloneMeta(meta), true
}

func Report(pass *analysis.Pass, opts ReportOptions) {
	pos := pass.Fset.Position(opts.ReportPos)
	evidenceFile := opts.EvidenceFile
	if evidenceFile == "" {
		evidenceFile = pos.Filename
	}

	meta := DiagnosticMeta{
		PkgPath:       pass.Pkg.Path(),
		Rule:          opts.Rule,
		Message:       opts.Message,
		ReportFile:    pos.Filename,
		ReportLine:    pos.Line,
		ReportColumn:  pos.Column,
		EvidenceFile:  evidenceFile,
		EvidenceLines: append([]int(nil), opts.EvidenceLines...),
		MatchMode:     opts.MatchMode,
	}

	Record(meta)
	pass.Report(analysis.Diagnostic{Pos: opts.ReportPos, Message: opts.Message})
}

func Reportf(pass *analysis.Pass, opts ReportOptions, format string, args ...interface{}) {
	opts.Message = fmt.Sprintf(format, args...)
	Report(pass, opts)
}

func cloneMeta(meta DiagnosticMeta) DiagnosticMeta {
	meta.EvidenceLines = append([]int(nil), meta.EvidenceLines...)
	return meta
}

func normalizeMeta(meta DiagnosticMeta) DiagnosticMeta {
	if meta.EvidenceFile == "" {
		meta.EvidenceFile = meta.ReportFile
	}
	if len(meta.EvidenceLines) == 0 && meta.ReportLine > 0 {
		meta.EvidenceLines = []int{meta.ReportLine}
	}
	return meta
}

func makeKey(pkgPath, file string, line, column int, message string) string {
	return fmt.Sprintf("%s|%s|%d|%d|%s", pkgPath, file, line, column, message)
}
