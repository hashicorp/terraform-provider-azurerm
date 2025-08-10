#!/bin/bash

#==============================================================================
# Cleanup Module
#==============================================================================

# Global cleanup flag
CLEANUP_MODE=false

# Main cleanup function
cleanup_ai_setup() {
    local repo_path="$1"
    
    log_info "Starting AI setup cleanup..."
    
    if [[ ! -d "$repo_path" ]]; then
        log_error "Repository path not found: $repo_path"
        return 1
    fi
    
    # Remove AI setup files
    remove_ai_files "$repo_path"
    
    # Restore backups
    restore_all_backups "$repo_path"
    
    # Clean up temporary files
    cleanup_temp_files
    
    log_success "AI setup cleanup completed!"
}

# Remove AI setup files
remove_ai_files() {
    local repo_path="$1"
    local ai_setup_dir="${repo_path}/.github/AILocalInstall"
    
    if [[ -d "$ai_setup_dir" ]]; then
        log_info "Removing AI setup directory..."
        rm -rf "$ai_setup_dir"
        log_success "AI setup directory removed"
    else
        log_warning "AI setup directory not found"
    fi
}

# Clean up temporary files
cleanup_temp_files() {
    log_info "Cleaning up temporary files..."
    
    local temp_dirs=(
        "/tmp/ai_local_install_$$"
        "/tmp/terraform_ai_setup_$$"
    )
    
    for temp_dir in "${temp_dirs[@]}"; do
        if [[ -d "$temp_dir" ]]; then
            rm -rf "$temp_dir"
            log_info "Removed temporary directory: $temp_dir"
        fi
    done
}

# Confirm cleanup operation
confirm_cleanup() {
    echo ""
    echo "⚠️  CLEANUP CONFIRMATION"
    echo "======================================"
    echo ""
    echo "This will:"
    echo "  • Remove all AI setup files"
    echo "  • Restore original configuration backups"
    echo "  • Clean up temporary files"
    echo ""
    
    if [[ "$AUTO_APPROVE" != "true" ]]; then
        read -p "Do you want to continue? (y/N): " -r
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Cleanup cancelled by user"
            return 1
        fi
    fi
    
    return 0
}

# Main cleanup entry point
perform_cleanup() {
    local repo_path="$1"
    
    log_info "=== CLEANUP MODE ==="
    
    if ! confirm_cleanup; then
        return 1
    fi
    
    if cleanup_ai_setup "$repo_path"; then
        echo "✅ Cleanup completed successfully!"
        return 0
    else
        log_error "Cleanup failed"
        return 1
    fi
}
