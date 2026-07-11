package list_upgrade

import (
	"os"
	"path/filepath"
	"strings"
)

// deriveReadModel resolves the SDK read model — the type of resp.Model returned
// by the Get method — by reading the vendored go-azure-sdk package. This is the
// authoritative source: the flatten method's model parameter and the list
// results slice must match resp.Model exactly. It avoids relying on the Pandora
// classification, which can be wrong for resources grouped into a single SDK
// package (e.g. everything under virtualwans).
func (r *Resource) deriveReadModel() {
	if r.GetMethod == "" || r.SDKImportPath == "" {
		return
	}
	root := findVendorRoot(r.Path)
	if root == "" {
		return
	}
	pkgDir := filepath.Join(root, "vendor", filepath.FromSlash(r.SDKImportPath))
	if m := readModelFromGetMethod(pkgDir, r.GetMethod); m != "" {
		r.ReadModel = m
	}
}

// findVendorRoot walks up from filePath to the directory containing a `vendor`
// folder (the provider module root), or "" when none is found.
func findVendorRoot(filePath string) string {
	dir, err := filepath.Abs(filePath)
	if err != nil {
		return ""
	}
	dir = filepath.Dir(dir)
	for {
		if info, err := os.Stat(filepath.Join(dir, "vendor")); err == nil && info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// readModelFromGetMethod finds `type <GetMethod>OperationResponse struct { ...
// Model *<X> ... }` in the vendored SDK package and returns X.
func readModelFromGetMethod(pkgDir, getMethod string) string {
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return ""
	}
	needle := "type " + getMethod + "OperationResponse struct"
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(pkgDir, e.Name()))
		if err != nil {
			continue
		}
		content := string(data)
		idx := strings.Index(content, needle)
		if idx == -1 {
			continue
		}
		if m := modelFieldType(content[idx:]); m != "" {
			return m
		}
	}
	return ""
}

// modelFieldType extracts the type X from the first `Model *X` field within the
// struct body that starts at the beginning of structText.
func modelFieldType(structText string) string {
	if end := strings.Index(structText, "}"); end != -1 {
		structText = structText[:end]
	}
	for _, line := range strings.Split(structText, "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 && fields[0] == "Model" {
			return strings.TrimPrefix(fields[1], "*")
		}
	}
	return ""
}
