package examplegen

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

const (
	testFileGen       = "example_gen_test.go"
	testCaseGenPrefix = "TestExampleGen_"
)

type ExampleSource struct {
	RootDir     string
	ServicePkgs []string
	TestCase    string
}

type TestFuncInfo struct {
	file     *ast.File
	filePath string
	fdecl    *ast.FuncDecl
}

type TestFuncInfos []TestFuncInfo

// TstFuncInfosSet stores the TestFuncInfos for each matched packages.
type TestFuncInfosSet map[*packages.Package]TestFuncInfos

func (src ExampleSource) GenExample() (string, error) {
	cfg := &packages.Config{Mode: packages.LoadAllSyntax, Tests: true, Dir: src.RootDir}
	pkgs, err := packages.Load(cfg, src.ServicePkgs...)
	if err != nil {
		return "", fmt.Errorf("loading package: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return "", fmt.Errorf("loading package contains error")
	}

	p := regexp.MustCompile(src.TestCase)

	testFuncInfosSet := TestFuncInfosSet{}
	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "_test") {
			continue
		}
		var testFuncInfos TestFuncInfos
		for _, f := range pkg.Syntax {
			for _, decl := range f.Decls {
				fdecl, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if !p.MatchString(fdecl.Name.Name) {
					continue
				}
				testFuncInfos = append(testFuncInfos, TestFuncInfo{
					file:     f,
					filePath: pkg.Fset.File(fdecl.Pos()).Name(),
					fdecl:    fdecl,
				})
			}
		}
		if len(testFuncInfos) == 0 {
			continue
		}
		testFuncInfosSet[pkg] = testFuncInfos
	}

	if len(testFuncInfosSet) == 0 {
		return "", fmt.Errorf("no test case found that matches %q", src.TestCase)
	}

	for pkg, testFuncInfos := range testFuncInfosSet {
		tc, err := testFuncInfos.GenPrintTestCase()
		if err != nil {
			return "", fmt.Errorf("generating the dummy test cases for %s: %w", pkg.PkgPath, err)
		}

		testFilePath := filepath.Join(filepath.Dir(testFuncInfos[0].filePath), testFileGen)
		content, err := imports.Process(testFilePath, []byte(tc), &imports.Options{Comments: false})
		if err != nil {
			return "", fmt.Errorf("imports processing the source code for %s: %v", testFilePath, err)
		}

		testFile, err := os.Create(testFilePath)
		if err != nil {
			return "", fmt.Errorf("creating test file %s: %w", testFilePath, err)
		}
		defer os.Remove(testFilePath)
		if _, err := testFile.Write(content); err != nil {
			return "", fmt.Errorf("writing to the test file %q: %v", testFilePath, err)
		}
	}

	// Run the test and fetch the printed Terraform configuration.
	args := []string{"test", "-v", "-run=" + testCaseGenPrefix}
	args = append(args, src.ServicePkgs...)
	cmd := exec.Command("go", args...)
	cmd.Dir = src.RootDir

	// The acceptance.BuildTestData depends on the following environment variables:
	// - ARM_TEST_LOCATION
	// - ARM_TEST_LOCATION_ALT1
	// - ARM_TEST_LOCATION_ALT2
	if os.Getenv("ARM_TEST_LOCATION") == "" {
		cmd.Env = append(os.Environ(), "ARM_TEST_LOCATION=WestEurope")
	}
	if os.Getenv("ARM_TEST_LOCATION_ALT1") == "" {
		cmd.Env = append(os.Environ(), "ARM_TEST_LOCATION_ALT1=WestUS")
	}
	if os.Getenv("ARM_TEST_LOCATION_ALT2") == "" {
		cmd.Env = append(os.Environ(), "ARM_TEST_LOCATION_ALT2=WestUS2")
	}

	var (
		stdout = bytes.NewBuffer([]byte{})
		stderr = bytes.NewBuffer([]byte{})
	)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("running the test: %v", err)
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			msg += "\n" + stderr.String()
		}
		return "", fmt.Errorf(msg)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	configSet := map[string]string{}
	var lines []string
	var name string
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "=== RUN"):
			name = strings.TrimSpace(strings.TrimPrefix(line, "=== RUN"))
		case strings.HasPrefix(line, "--- PASS"):
			configSet[name] = string(hclwrite.Format([]byte(strings.Join(lines, "\n"))))
			lines = []string{}
			name = ""
		case strings.HasPrefix(line, "--- FAIL"):
			configSet[name] = fmt.Sprintf("Runtime Error:\n%s", strings.Join(lines, "\n"))
			lines = []string{}
			name = ""
		default:
			if name != "" && line != "" {
				lines = append(lines, line)
			}
		}
	}

	names := []string{}
	for k := range configSet {
		names = append(names, k)
	}
	sort.Strings(names)
	var out string
	for _, name := range names {
		out += fmt.Sprintf("# Generated from AccTest %s\n\n%s\n\n", strings.TrimPrefix(name, testCaseGenPrefix), configSet[name])
	}
	return out, nil
}

func (infos TestFuncInfos) GenPrintTestCase() (string, error) {
	if len(infos) == 0 {
		return "", nil
	}

	decls := []ast.Decl{}
	for _, info := range infos {
		fdecl, err := info.GenPrintTestCaseFdecl()
		if err == nil {
			decls = append(decls, fdecl)
		} else {
			decls = append(decls, info.GenParseErrorTestCase(info.fdecl.Name.Name, err))
		}
	}

	fset := token.NewFileSet()
	file := *infos[0].file
	file.Decls = decls

	buf := bytes.NewBuffer([]byte{})
	if err := printer.Fprint(buf, fset, &file); err != nil {
		return "", fmt.Errorf("printing the source code: %v", err)
	}
	return buf.String(), nil
}

func (info TestFuncInfo) GenPrintTestCaseFdecl() (*ast.FuncDecl, error) {
	// We keep the content of the target function body until the invocation of the "ResourceTest"/"DataSourceTest", which is typically doing the setup.
	var (
		resourceTestCall   *ast.CallExpr
		dataSourceTestCall *ast.CallExpr
	)
	var stmts []ast.Stmt
	for _, stmt := range info.fdecl.Body.List {
		exprstmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			stmts = append(stmts, stmt)
			continue
		}
		callexpr, ok := exprstmt.X.(*ast.CallExpr)
		if !ok {
			stmts = append(stmts, stmt)
			continue
		}
		selexpr, ok := callexpr.Fun.(*ast.SelectorExpr)
		if !ok {
			stmts = append(stmts, stmt)
			continue
		}
		if selexpr.Sel.Name == "ResourceTest" {
			resourceTestCall = callexpr
			break
		}
		if selexpr.Sel.Name == "DataSourceTest" {
			dataSourceTestCall = callexpr
			break
		}
		stmts = append(stmts, stmt)
		continue
	}

	if resourceTestCall == nil && dataSourceTestCall == nil {
		return nil, fmt.Errorf("no ResourceTest/DataSourceTest call found")
	}

	// Look for the first test step and record the assignment to the "Config", which is then used as the config generation invocation.
	var teststeps *ast.CompositeLit

	if resourceTestCall != nil {
		if len(resourceTestCall.Args) != 3 {
			return nil, fmt.Errorf("ResourceTest doesn't have 3 arguments")
		}
		var ok bool
		teststeps, ok = resourceTestCall.Args[2].(*ast.CompositeLit)
		if !ok {
			return nil, fmt.Errorf("test steps are not defined as a composite literal")
		}
	} else {
		if len(dataSourceTestCall.Args) != 2 {
			return nil, fmt.Errorf("DataSourceTest doesn't have 2 arguments")
		}
		var ok bool
		teststeps, ok = dataSourceTestCall.Args[1].(*ast.CompositeLit)
		if !ok {
			return nil, fmt.Errorf("test steps are not defined as a composite literal")
		}
	}
	if len(teststeps.Elts) == 0 {
		return nil, fmt.Errorf("there is no test step defined")
	}

	// The first test step can be either:
	// - A resource.TestStep composite literal: e.g. https://github.com/hashicorp/terraform-provider-azurerm/blob/fcfcdbdf8051050f83c0ec39d09b2bc2ca9152a9/internal/services/network/subnet_resource_test.go#L22
	// - A call to the data.ApplyStep: e.g. https://github.com/hashicorp/terraform-provider-azurerm/blob/fcfcdbdf8051050f83c0ec39d09b2bc2ca9152a9/internal/services/resource/resource_group_resource_test.go#L24
	var configcallexpr *ast.CallExpr
	switch firstteststep := teststeps.Elts[0].(type) {
	case *ast.CompositeLit:
		for _, elt := range firstteststep.Elts {
			elt, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := elt.Key.(*ast.Ident)
			if !ok {
				continue
			}
			if key.Name != "Config" {
				continue
			}
			callexpr, ok := elt.Value.(*ast.CallExpr)
			if !ok {
				continue
			}
			configcallexpr = callexpr
			break
		}
	case *ast.CallExpr:
		applystep, ok := firstteststep.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil, fmt.Errorf("the first test step is not a call of selector expression")
		}
		if applystep.Sel.Name != "ApplyStep" {
			return nil, fmt.Errorf("the first test step is not a call of data.ApplyStep")
		}
		configcallexpr = &ast.CallExpr{
			Fun:  firstteststep.Args[0],
			Args: []ast.Expr{applystep.X},
		}
	}

	if configcallexpr == nil {
		return nil, fmt.Errorf(`no "Config" field found for the first test step`)
	}

	stmts = append(stmts,
		&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "fmt",
					},
					Sel: &ast.Ident{
						Name: "Println",
					},
				},
				Args: []ast.Expr{
					configcallexpr,
				},
			},
		},
	)

	fdecl := *info.fdecl
	body := *fdecl.Body
	body.List = stmts
	fdecl.Body = &body
	name := *fdecl.Name
	name.Name = info.TestCaseName()
	fdecl.Name = &name

	return &fdecl, nil
}

func (info TestFuncInfo) TestCaseName() string {
	return fmt.Sprintf("%s%s", testCaseGenPrefix, info.fdecl.Name.Name)
}

func (info TestFuncInfo) GenParseErrorTestCase(pkg string, err error) *ast.FuncDecl {
	src := fmt.Sprintf(`
package %s

func %s(t *testing.T) {
	fmt.Println(%q)
}
`, pkg, info.TestCaseName(), "Parsing Error: "+err.Error())
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "src.go", src, 0)
	return f.Decls[0].(*ast.FuncDecl)
}
