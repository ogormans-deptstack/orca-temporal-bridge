# Branch Protection Rules — Main Branch

**Status**: Ready to enable before sharing with colleagues

**Policy**: Automated safety gates on every merge to main

## Rules to Enable

### 1. Require Status Checks to Pass (Before Merge)

```
✅ verify-versions
✅ test-local
✅ test-coverage
✅ test-bridge
✅ test-events
✅ test-integration
✅ build
✅ lint (optional, can be slow)
```

**Why**: Ensures every merge has passed all CI checks. No broken main.

### 2. Require Branches to Be Up-to-Date Before Merging

```
✓ Enabled
```

**Why**: If main changes while your PR is open, GitHub blocks merge until rebased. Prevents conflicts.

### 3. Require Pull Request Reviews Before Merging

```
Number of required reviews: 1
Dismiss stale pull request approvals when new commits are pushed: ✓
Require code owner review: ✗ (optional, enable if CODEOWNERS exists)
```

**Why**: At least one human review before merging. Stale reviews are dismissed so new commits get fresh eyes.

### 4. Require Conversation Resolution

```
✓ Enabled
```

**Why**: All review comments must be resolved. No hidden blockers.

### 5. Require Signed Commits (Optional but Recommended)

```
✓ Enabled for team security
```

**Why**: Every commit must be GPG-signed. Verifies identity.

### 6. Restrict Pushes (Who Can Push)

```
Allow specified actors to bypass branch protections:
  - (none initially, can be configured later)
```

**Why**: No direct pushes to main. Everything goes through PR + reviews.

### 7. Require Linear History (Optional)

```
✓ Enabled
```

**Why**: Keeps git history clean. Squash or rebase, no merge commits.

## How to Apply (GitHub UI)

### Step 1: Navigate to Branch Protection
1. Go to **Settings** → **Branches** → **Branch protection rules**
2. Click **Add rule**
3. Pattern: `main`

### Step 2: Configure Each Rule

| Rule | Setting |
|------|---------|
| **Protect matching branches** | ✓ Checked |
| **Require status checks to pass before merging** | ✓ Checked |
| **Require branches to be up to date** | ✓ Checked |
| **Require pull request reviews** | ✓ Checked (1 required) |
| **Dismiss stale reviews** | ✓ Checked |
| **Require review from code owners** | ✗ Unchecked (unless CODEOWNERS exists) |
| **Require conversation resolution** | ✓ Checked |
| **Require signed commits** | ✓ Checked |
| **Require deployment to succeed** | ✗ (not yet, enable in Sprint 1.4) |
| **Lock branch** | ✗ (not needed) |
| **Allow force pushes** | ✗ Never allow |
| **Allow deletions** | ✗ Never allow |

### Step 3: Save

Click **Create** or **Update** rule.

## How to Apply (GitHub CLI)

```bash
# Install GitHub CLI if not already installed
brew install gh

# Authenticate
gh auth login

# Create branch protection rule
gh api repos/ogormans-deptstack/orca-temporal-bridge/branches/main/protection \
  -X PUT \
  -f required_status_checks='{"strict":true,"contexts":["verify-versions","test-local","test-coverage","test-bridge","test-events","test-integration","build"]}' \
  -f required_pull_request_reviews='{"dismissal_restrictions":{},"dismiss_stale_reviews":true,"require_code_owner_reviews":false,"required_approving_review_count":1}' \
  -f enforce_admins=false \
  -f require_linear_history=true \
  -f require_conversation_resolution=true \
  -f require_signed_commits=true \
  -f allow_force_pushes=false \
  -f allow_deletions=false

# Verify the rule was created
gh api repos/ogormans-deptstack/orca-temporal-bridge/branches/main/protection
```

## Testing the Protection

### Test 1: Direct Push Should Fail
```bash
git push origin main  # Should be rejected
# Error: "Protected branch" / "Rules preventing push"
```

### Test 2: PR Without Reviews Should Not Merge
```bash
# Create a PR, click "Merge" without approvals
# Should show: "1 approval required" — blocked
```

### Test 3: Failed CI Should Block Merge
```bash
# Edit .tool-versions to break Go version
# Create a PR
# CI fails (verify-versions job)
# Merge button disabled until CI passes
```

### Test 4: Stale Reviews Should Be Dismissed
```bash
# Get 1 approval on a PR
# Push a new commit to the PR
# Approval count resets to 0 (requires new review)
```

## For Sharing with Colleagues

Once protection is enabled:

1. **New contributors**: Explain "We use main branch protection — all changes go through PR"
2. **PR workflow**: "Create PR → Wait for CI (2 min) → Get 1 review → Merge"
3. **Emergency fixes**: "If urgent, create PR and ping a maintainer for quick review"
4. **CI failures**: "Check the sprint summary report artifact to debug"

## Future Enhancements (Sprint 1.4+)

- [ ] Add code owner review requirement (create CODEOWNERS file)
- [ ] Add deployment requirement (after dev/staging env setup)
- [ ] Add automatic label rules (e.g., `bug` PRs need security review)
- [ ] Add automatic issue links (e.g., `Closes #123` requirement)
- [ ] Add Slack notification on protection violations

## Disable/Modify Later

If protection is too strict:
1. Go to **Settings** → **Branches**
2. Click the rule
3. Adjust settings
4. Or delete the rule entirely

---

**Status**: Ready to enable. Test locally first (Steps 1-4 above), then enable for colleagues.
