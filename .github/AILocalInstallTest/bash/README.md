# Bash Test Suite

Test suite for the bash installer (`install-copilot-setup.sh`).

## Test Files

### `test_installer.sh`
Comprehensive test suite for bash installer functionality.

**What it tests:**
- Help functionality and argument parsing
- Prerequisites checking (bash version, git, etc.)
- Repository validation and detection
- Module loading and sourcing
- Color output functions
- Error handling and edge cases
- User interface components

## Usage

```bash
# From the test directory
cd .github/AILocalInstallTest/bash
chmod +x test_installer.sh
./test_installer.sh
```

## Requirements

- Bash 3.2+
- Git (for repository detection tests)
- Access to `../../AILocalInstall/` directory
- Standard Unix utilities (grep, awk, etc.)

## Test Results

- ✅ **PASS**: Component works correctly
- ❌ **FAIL**: Component has issues  
- ℹ️ **INFO**: Additional information

## Adding Tests

To add new bash-specific tests:

1. Add test logic to `test_installer.sh`
2. Follow the existing pattern:
   ```bash
   echo "✅ Test N: Description"
   # Test logic here
   if [condition]; then
       echo "   PASS: Test description"
   else
       echo "   FAIL: Test description" 
   fi
   ```
3. Update this README with test description
