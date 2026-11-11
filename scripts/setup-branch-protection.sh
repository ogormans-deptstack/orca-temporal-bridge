#!/bin/bash
#
# Setup Branch Protection for Main Branch
# 
# This script configures GitHub branch protection rules for the main branch.
# Requires: GitHub CLI (gh), authenticated with proper permissions
#
# Usage:
#   bash scripts/setup-branch-protection.sh
#   
# The script will:
#   ✓ Verify GitHub CLI is installed
#   ✓ Check repository permissions
#   ✓ Apply branch protection rules
#   ✓ Verify the rules were applied correctly
#

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_OWNER="ogormans-deptstack"
REPO_NAME="orca-temporal-bridge"
BRANCH="main"

# Functions
print_info() { echo -e "${BLUE}ℹ${NC} $1"; }
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_warning() { echo -e "${YELLOW}⚠${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }

# Verify requirements
verify_requirements() {
    print_info "Verifying requirements..."
    
    if ! command -v gh &> /dev/null; then
        print_error "GitHub CLI not found. Install: brew install gh"
        exit 1
    fi
    print_success "GitHub CLI installed: $(gh --version)"
    
    # Check authentication
    if ! gh auth status &> /dev/null; then
        print_error "Not authenticated with GitHub. Run: gh auth login"
        exit 1
    fi
    print_success "GitHub authentication verified"
    
    # Check repository access
    if ! gh repo view "${REPO_OWNER}/${REPO_NAME}" &> /dev/null; then
        print_error "Repository not accessible. Check owner/name: ${REPO_OWNER}/${REPO_NAME}"
        exit 1
    fi
    print_success "Repository accessible: ${REPO_OWNER}/${REPO_NAME}"
}

# Apply branch protection rules
apply_protection() {
    print_info ""
    print_info "Applying branch protection rules to '${BRANCH}'..."
    echo ""
    
    # Define the protection rules as JSON
    PROTECTION_RULES=$(cat <<'EOF'
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "verify-versions",
      "test-local",
      "test-coverage",
      "test-bridge",
      "test-events",
      "test-integration",
      "build"
    ]
  },
  "required_pull_request_reviews": {
    "dismissal_restrictions": {},
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": false,
    "required_approving_review_count": 1
  },
  "require_code_owner_reviews": false,
  "require_last_push_approval": false,
  "enforce_admins": false,
  "require_linear_history": true,
  "require_conversation_resolution": true,
  "require_signed_commits": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "restrictions": null
}
EOF
)
    
    # Apply the rules via GitHub API
    if gh api \
        --method PUT \
        "repos/${REPO_OWNER}/${REPO_NAME}/branches/${BRANCH}/protection" \
        --input <(echo "$PROTECTION_RULES"); then
        
        print_success "Branch protection rules applied successfully"
    else
        print_error "Failed to apply branch protection rules"
        print_warning "Make sure you have admin permissions on the repository"
        exit 1
    fi
}

# Verify protection was applied
verify_protection() {
    print_info ""
    print_info "Verifying branch protection rules..."
    echo ""
    
    # Fetch and display current rules
    CURRENT_RULES=$(gh api "repos/${REPO_OWNER}/${REPO_NAME}/branches/${BRANCH}/protection")
    
    print_success "Branch protection is enabled for '${BRANCH}'"
    echo ""
    echo "Current rules:"
    echo "$CURRENT_RULES" | jq '.'
}

# Display summary
print_summary() {
    echo ""
    print_info "========================================="
    print_info "Branch Protection Summary"
    print_info "========================================="
    echo ""
    
    echo "Repository: ${REPO_OWNER}/${REPO_NAME}"
    echo "Branch: ${BRANCH}"
    echo ""
    
    echo "Rules applied:"
    echo "  ✓ Require status checks: verify-versions, test-*, build"
    echo "  ✓ Require 1 pull request review"
    echo "  ✓ Dismiss stale reviews when new commits are pushed"
    echo "  ✓ Require conversation resolution"
    echo "  ✓ Require signed commits"
    echo "  ✓ Require branches to be up-to-date"
    echo "  ✓ Require linear history (squash or rebase only)"
    echo "  ✗ No direct pushes allowed"
    echo "  ✗ No force pushes"
    echo ""
    
    echo "Testing the protection:"
    echo "  1. Try to push directly: git push origin main"
    echo "     Expected: 'remote: error: protected branch'"
    echo ""
    echo "  2. Create a PR without approval and try to merge"
    echo "     Expected: 'Merging is blocked' (needs 1 approval)"
    echo ""
    echo "  3. Approve and merge"
    echo "     Expected: Success (all checks passed, 1 approval, etc.)"
    echo ""
    
    print_success "Branch protection is now active!"
    echo ""
    print_info "For details, see: BRANCH_PROTECTION.md"
}

# Main flow
main() {
    echo ""
    echo "╔════════════════════════════════════════════════════╗"
    echo "║   Branch Protection Setup — ORCA-Temporal Bridge    ║"
    echo "╚════════════════════════════════════════════════════╝"
    echo ""
    
    verify_requirements
    echo ""
    
    print_warning "This will enable branch protection on '${BRANCH}'"
    print_warning "Once enabled, direct pushes and unreviewed PRs cannot merge."
    echo ""
    
    read -p "Continue? (y/n): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Cancelled."
        exit 0
    fi
    
    apply_protection
    verify_protection
    print_summary
}

# Run
main
