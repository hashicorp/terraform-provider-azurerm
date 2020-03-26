package analysisfixtest

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/singlechecker"
)

const (
	// Directory name suffix used for expected fix files
	FixedDirectorySuffix = `_fixed`
)

// Run ensures that Analyzer SuggestedFixes produce the expected file contents
// This is done by copying source testdata files into a temporary directory,
// running the Analyzer with the -fix flag, and comparing the contents of every
// expected file to the fixed file contents.
func Run(t *testing.T, testdataDir string, analyzer *analysis.Analyzer, packageName string) {
	// The go/analysis framework hides the functionality to trigger suggested fixes
	// in internal packages, so we wrap the execution of the Analyzer in an
	// exported Main function for testing. For reference:
	// https://github.com/golang/tools/blob/521f4a0cd458c441dec3bafd8ba24526a6cb9b09/go/analysis/multichecker/multichecker_test.go
	if os.Getenv("TEST_CHILD") == "1" {
		// replace [progname -test.run=TestName -- ...]
		//      by [progname ...]
		os.Args = os.Args[2:]
		os.Args[0] = "vet"
		singlechecker.Main(analyzer)
		panic("unreachable")
	}

	// Source files to copy into temporary directory
	srcDir := filepath.Join(testdataDir, "src", packageName)
	// Expected files to compare after fixes
	wantDir := filepath.Join(testdataDir, "src", packageName+FixedDirectorySuffix)

	srcFiles, err := readFiles(srcDir)

	if err != nil {
		t.Fatal(err)
	}

	tmpDir, tmpDirCleanup, err := analysistest.WriteFiles(srcFiles)

	if err != nil {
		t.Fatalf("error preparing test files: %s", err)
	}

	defer tmpDirCleanup()

	gotDir := filepath.Join(tmpDir, "src", packageName)

	if err := addVendor(srcDir, gotDir); err != nil {
		t.Fatal(err)
	}

	runAnalyzerFixes(t, gotDir)

	compareFiles(t, wantDir, gotDir)
}

func addVendor(srcdir string, dstdir string) error {
	src := filepath.Join(srcdir, "vendor")
	fi, err := os.Stat(src)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading (%s): %w", src, err)
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("only vendor symlink handling implemented")
	}

	vendorDst, err := filepath.EvalSymlinks(src)

	if err != nil {
		return fmt.Errorf("error evaluating symlink (%s): %w", src, err)
	}

	dst := filepath.Join(dstdir, "vendor")

	if err := os.Symlink(vendorDst, dst); err != nil {
		return fmt.Errorf("error creating (%s) symlink (%s): %w", vendorDst, dst, err)
	}

	return nil
}

func compareFiles(t *testing.T, wantDir string, gotDir string) {
	gotFiles, err := readFiles(gotDir)

	if err != nil {
		t.Fatal(err)
	}

	wantFiles, err := readFiles(wantDir)

	if err != nil {
		t.Fatal(err)
	}

	dmp := diffmatchpatch.New()

	for wantPath, wantContents := range wantFiles {
		wantPath = strings.Replace(wantPath, FixedDirectorySuffix, "", 1)
		gotContents, ok := gotFiles[wantPath]

		if !ok {
			t.Errorf("missing updated file (%s) content", wantPath)
		}

		if wantContents == gotContents {
			continue
		}

		wantRunes, gotRunes, lines := dmp.DiffLinesToRunes(wantContents, gotContents)
		diffs := dmp.DiffMainRunes(wantRunes, gotRunes, false)
		diffs = dmp.DiffCharsToLines(diffs, lines)
		t.Errorf("unexpected differences in file (%s) contents: \n%s\n", wantPath, dmp.DiffPrettyText(diffs))
	}
}

func readFiles(dir string) (map[string]string, error) {
	dirEntries, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, fmt.Errorf("error reading directory (%s): %w", dir, err)
	}

	var fullpaths []string

	for _, dirEntry := range dirEntries {
		// Skip symlinks as we will handle vendor separately
		if dirEntry.Mode()&os.ModeSymlink != 0 {
			continue
		}

		if dirEntry.IsDir() {
			// Only reading one level deep is okay for our purposes
			subdir := filepath.Join(dir, dirEntry.Name())
			newfullpaths, err := filepath.Glob(filepath.Join(subdir, "*.go"))

			if err != nil {
				return nil, fmt.Errorf("error reading directory (%s) files: %w", subdir, err)
			}

			fullpaths = append(fullpaths, newfullpaths...)
			continue
		}

		if !strings.HasSuffix(dirEntry.Name(), ".go") {
			continue
		}

		fullpaths = append(fullpaths, filepath.Join(dir, dirEntry.Name()))
	}

	filemap := make(map[string]string)

	for _, fullpath := range fullpaths {
		contents, err := ioutil.ReadFile(fullpath)

		if err != nil {
			return nil, fmt.Errorf("error reading file (%s): %w", fullpath, err)
		}

		relpath, err := filepath.Rel(filepath.Dir(dir), fullpath)

		if err != nil {
			return nil, fmt.Errorf("error getting relative path between (%s, %s): %w", filepath.Dir(dir), fullpath, err)
		}

		filemap[relpath] = string(contents)
	}

	return filemap, nil
}

func runAnalyzerFixes(t *testing.T, gotDir string) {
	args := []string{
		fmt.Sprintf("-test.run=%s", t.Name()),
		"--",
		"-fix",
		gotDir,
	}
	envs := []string{
		fmt.Sprintf("GOPATH=%s", filepath.Dir(filepath.Dir(gotDir))), // trim src/PACKAGE directories
		"GO111MODULE=off",
		"TEST_CHILD=1",
	}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), envs...)
	out, err := cmd.CombinedOutput()

	if len(out) > 0 {
		t.Logf("%s: out=<<%s>>", args, out)
	}

	var exitcode int

	if err, ok := err.(*exec.ExitError); ok {
		exitcode = err.ExitCode()
	}

	// Diagnostics return exit code 3
	if exitcode != 3 {
		t.Errorf("%s: exited %d, want %d", args, exitcode, 3)
	}
}
