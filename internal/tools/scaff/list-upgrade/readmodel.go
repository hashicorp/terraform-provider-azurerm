// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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

// deriveListMethods resolves the list operations for a non-parent-scoped
// resource by scanning the vendored SDK for the list Complete methods, keyed by
// their id parameter type. A SubscriptionId / ResourceGroupId parameter yields a
// top-level list; any other resource id (e.g. commonids.VirtualNetworkId for
// subnets.ListComplete) identifies a parent-scoped list, from which the parent
// scope is derived when the resource source did not reveal a parent. This
// mirrors how the read model is derived, keeping the generated list ops in sync
// with the SDK rather than relying on Pandora's classification.
func (r *Resource) deriveListMethods() {
	if r.ParentIDType != "" || r.GetMethod == "" || r.SDKImportPath == "" {
		return
	}
	root := findVendorRoot(r.Path)
	if root == "" {
		return
	}
	pkgDir := filepath.Join(root, "vendor", filepath.FromSlash(r.SDKImportPath))
	prefix := strings.TrimSuffix(r.GetMethod, "Get")
	m := listMethodsFromVendor(pkgDir, prefix)
	r.ListSubscriptionMethod = m.sub
	r.ListResourceGroupMethod = m.rg
	// When the SDK exposes no subscription/resource-group list but does expose a
	// parent-scoped list, derive the parent scope from the list method's id
	// parameter type (e.g. VirtualNetworkId -> virtual_network_id).
	if m.sub == "" && m.rg == "" && m.parent != "" && m.parentIDType != "" {
		r.setSDKParent(m.parent, m.parentIDType, m.parentIDPkg)
	}
}

// vendorListMethods holds the list methods discovered in a vendored SDK package
// for a given prefix, classified by their id parameter type.
type vendorListMethods struct {
	sub          string // subscription-scoped list method (name without Complete)
	rg           string // resource-group-scoped list method
	parent       string // parent-scoped list method
	parentIDType string // the parent id type, e.g. "VirtualNetworkId"
	parentIDPkg  string // the parent id package, e.g. "commonids"
}

// listMethodsFromVendor scans the vendored SDK package for `<name>Complete`
// methods (name starting with prefix and containing "List") and classifies them
// by their id parameter type.
func listMethodsFromVendor(pkgDir, prefix string) vendorListMethods {
	var m vendorListMethods
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return m
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(pkgDir, e.Name()))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			name, idType, idPkg := parseListCompleteMethod(line, prefix)
			if name == "" {
				continue
			}
			switch idType {
			case "SubscriptionId":
				if m.sub == "" {
					m.sub = name
				}
			case "ResourceGroupId":
				if m.rg == "" {
					m.rg = name
				}
			default:
				if m.parent == "" {
					m.parent = name
					m.parentIDType = idType
					m.parentIDPkg = idPkg
				}
			}
		}
	}
	return m
}

// parseListCompleteMethod parses a `func (c XClient) <name>Complete(ctx
// context.Context, id <pkg>.<idType>)` declaration, returning the method name
// without the Complete suffix, the id type and its package, when name starts
// with prefix and contains "List".
func parseListCompleteMethod(line, prefix string) (name, idType, idPkg string) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "func (c ") {
		return "", "", ""
	}
	rp := strings.Index(line, ") ")
	if rp == -1 {
		return "", "", ""
	}
	rest := line[rp+2:]
	const sig = "Complete(ctx context.Context, id "
	cp := strings.Index(rest, sig)
	if cp == -1 {
		return "", "", ""
	}
	name = rest[:cp]
	if !strings.HasPrefix(name, prefix) || !strings.Contains(name, "List") {
		return "", "", ""
	}
	after := rest[cp+len(sig):]
	end := strings.IndexAny(after, "),")
	if end <= 0 {
		return "", "", ""
	}
	typeExpr := strings.TrimSpace(after[:end])
	if typeExpr == "" {
		return "", "", ""
	}
	if dot := strings.LastIndex(typeExpr, "."); dot != -1 {
		return name, typeExpr[dot+1:], typeExpr[:dot]
	}
	return name, typeExpr, ""
}
