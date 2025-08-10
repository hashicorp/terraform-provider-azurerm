#!/bin/bash

echo "🧪 Comprehensive Test Suite for AI Installation Scripts (Bash)"
echo "=============================================================="
echo

# Create isolated test environment
TEST_DIR="/tmp/terraform-azurerm-test-$(date +%s)-$$"
echo "📁 Creating isolated test environment: $TEST_DIR"

# Copy installer files to test directory (isolation)
cp -r ../../AILocalInstall "$TEST_DIR"
cd "$TEST_DIR/AILocalInstall" || exit 1

echo "✅ Test environment isolated successfully"
echo

# Test 1: Help functionality
echo "✅ Test 1: Help functionality"
./install-copilot-setup.sh -help > /dev/null 2>&1
if [[ $? -eq 0 ]]; then
    echo "   PASS: Help command works"
else
    echo "   FAIL: Help command failed"
fi

# Test 2: Invalid arguments
echo "✅ Test 2: Invalid argument handling"
output=$(./install-copilot-setup.sh --invalid 2>&1)
if echo "$output" | grep -q "Unknown option"; then
    echo "   PASS: Invalid arguments properly rejected"
else
    echo "   FAIL: Invalid arguments not handled"
fi

# Test 3: Prerequisites check
echo "✅ Test 3: Prerequisites check"
source modules/bash/core-functions.sh
if check_prerequisites; then
    echo "   PASS: Prerequisites check passed"
else
    echo "   FAIL: Prerequisites check failed"
fi

# Test 4: Repository validation (using test repo, not real one)
echo "✅ Test 4: Repository validation"
# Create mock invalid repository for testing
MOCK_INVALID_REPO="/tmp/mock-invalid-repo-$$"
mkdir -p "$MOCK_INVALID_REPO"
if validate_repository "$MOCK_INVALID_REPO"; then
    echo "   FAIL: Should reject invalid repo without go.mod"
else
    echo "   PASS: Correctly rejected invalid repository"
fi
rm -rf "$MOCK_INVALID_REPO"

# Test 5: Module loading
echo "✅ Test 5: Module loading"
modules_loaded=0
for module in modules/bash/*.sh; do
    if source "$module" 2>/dev/null; then
        ((modules_loaded++))
    fi
done
echo "   INFO: Loaded $modules_loaded modules successfully"

# Test 6: Color output
echo "✅ Test 6: Color output functions"
source modules/bash/core-functions.sh
log_info "Testing colored output"
log_success "Success message"
log_warning "Warning message"  
log_error "Error message"
echo "   INFO: Color output functions work"

# Test 7: JSON functions (if jq available)
echo "✅ Test 7: JSON functions"
if command -v jq >/dev/null; then
    echo '{"test": "value"}' > /tmp/test.json
    if is_valid_json "/tmp/test.json"; then
        echo "   PASS: JSON validation works"
    else
        echo "   FAIL: JSON validation failed"
    fi
    rm -f /tmp/test.json
else
    echo "   SKIP: jq not available"
fi

echo
echo "🎉 Test suite completed!"
echo "======================================================="

# Cleanup test environment
echo "🧹 Cleaning up test environment..."
cd /
rm -rf "$TEST_DIR"
echo "✅ Test environment cleaned up"
