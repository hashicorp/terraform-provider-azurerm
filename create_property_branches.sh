#!/bin/bash

# Script to create separate branches for each AKS node pool property

# Properties to implement
PROPERTIES=(
    "gateway-profile"
    "pod-ip-allocation-mode" 
    "security-profile"
    "virtual-machine-nodes-status"
    "virtual-machine-profile"
)

# Base branch name
BASE_BRANCH="feature/aks-node-pool"

echo "Creating branches for AKS node pool properties..."

for property in "${PROPERTIES[@]}"; do
    branch_name="${BASE_BRANCH}-${property}"
    echo "Creating branch: $branch_name"
    
    # Create and checkout new branch
    git checkout -b "$branch_name"
    
    # Reset to original state (before our changes)
    git reset --hard HEAD~1
    
    echo "Branch $branch_name created. Ready to implement $property property."
    echo "Next steps:"
    echo "1. Implement only the $property property"
    echo "2. Add tests for $property"
    echo "3. Update documentation for $property"
    echo "4. Commit changes"
    echo ""
    
    # Go back to main for next iteration
    git checkout main
done

echo "All property branches created successfully!"
echo "Each branch is ready for individual property implementation."
