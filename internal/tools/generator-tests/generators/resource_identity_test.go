package generators

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestResourceIdentity_ParentResource(t *testing.T) {
	tmpDir := t.TempDir()

	// Create refactored subpackage structure:
	// internal/services/network/resources/network_watcher/
	// internal/services/network/resources/network_watcher_flow_log/
	baseDir := filepath.Join(tmpDir, "internal", "services", "network", "resources")
	watcherDir := filepath.Join(baseDir, "network_watcher")
	flowLogDir := filepath.Join(baseDir, "network_watcher_flow_log")

	if err := os.MkdirAll(watcherDir, 0o755); err != nil {
		t.Fatalf("failed to create watcherDir: %v", err)
	}
	if err := os.MkdirAll(flowLogDir, 0o755); err != nil {
		t.Fatalf("failed to create flowLogDir: %v", err)
	}

	// Create a minimal resource.go in child dir so isRefactoredSubpackage recognizes it
	flowLogResourceGo := filepath.Join(flowLogDir, "resource.go")
	if err := os.WriteFile(flowLogResourceGo, []byte("package network_watcher_flow_log\n"), 0o644); err != nil {
		t.Fatalf("failed to write dummy resource.go: %v", err)
	}

	// Change cwd to child directory
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	if err := os.Chdir(flowLogDir); err != nil {
		t.Fatalf("failed to chdir to flowLogDir: %v", err)
	}

	cmd := &ResourceIdentityCommand{
		Ui: &cli.MockUi{},
	}

	args := []string{
		"-resource-name", "network_watcher_flow_log",
		"-properties", "name,resource_group_name,network_watcher_name",
		"-test-sequential",
		"-parent-resource", "network_watcher",
	}

	if code := cmd.Run(args); code != 0 {
		mockUi, _ := cmd.Ui.(*cli.MockUi)
		t.Fatalf("cmd.Run failed with exit code %d: %s", code, mockUi.ErrorWriter.String())
	}

	// Expected output: flow_log_resource_identity_gen_test.go in watcherDir
	expectedPath := filepath.Join(watcherDir, "flow_log_resource_identity_gen_test.go")
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("expected generated file at %s: %v", expectedPath, err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "package network_watcher_test") {
		t.Errorf("expected package network_watcher_test, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "func testAccNetworkWatcherFlowLog_resourceIdentity(t *testing.T)") {
		t.Errorf("expected func testAccNetworkWatcherFlowLog_resourceIdentity, got:\n%s", contentStr)
	}

	// Verify no file was left in flowLogDir
	stalePath := filepath.Join(flowLogDir, "resource_identity_gen_test.go")
	if _, err := os.Stat(stalePath); err == nil {
		t.Errorf("stale file should not exist at %s", stalePath)
	}
}
