# Sprint 1.4 Checklist: Docker Login & Lint Fixes

**Sprint:** 1.4  
**Objective:** Enable Docker Temporal integration and code style strictness  
**Tracking:** Use this checklist for implementation  

---

## Task 1: Docker Hub/GitHub Login for Temporal 🐳

### Step 1: Get Credentials
- [ ] Create GitHub PAT (Personal Access Token)
  - Go to: https://github.com/settings/tokens
  - Select "Tokens (classic)"
  - Scopes: `read:packages`
  - Copy token (only shown once!)
  
- [ ] OR use Docker Hub credentials
  - Username: your Docker Hub login
  - Token: Docker Hub personal access token (Settings → Security → Personal Access Tokens)

### Step 2: Add GitHub Secrets
- [ ] In repo settings: https://github.com/ogormans-deptstack/orca-temporal-bridge/settings/secrets/actions
  - [ ] Add secret: `DOCKERHUB_USERNAME` (GitHub username or Docker Hub username)
  - [ ] Add secret: `DOCKERHUB_TOKEN` (your PAT token)

### Step 3: Update CI Workflow
- [ ] Edit `.github/workflows/test.yml`
- [ ] Find: `Integration Tests (with Temporal)` job
- [ ] Add docker login step (AFTER checkout, BEFORE go setup):
  ```yaml
  - name: Log in to Docker
    uses: docker/login-action@v2
    with:
      username: ${{ secrets.DOCKERHUB_USERNAME }}
      password: ${{ secrets.DOCKERHUB_TOKEN }}
  ```
  
- [ ] OR for GitHub Container Registry:
  ```yaml
  - name: Log in to GitHub Container Registry
    uses: docker/login-action@v2
    with:
      registry: ghcr.io
      username: ${{ github.actor }}
      password: ${{ secrets.GITHUB_TOKEN }}
  ```

### Step 4: Test Locally
```bash
# Simulate docker pull
docker pull temporalapi/server:latest
# or: docker pull ghcr.io/temporalapi/server:latest

# Run full test
docker-compose up -d
make test-with-temporal
docker-compose down
```

### Step 5: Verify in CI
- [ ] Push branch with docker login changes
- [ ] Check GitHub Actions run
- [ ] Verify: `Integration Tests (with Temporal)` job passes ✅
- [ ] No more "pull access denied" errors

**Definition of Done:** Integration tests pass in CI for 2 consecutive runs

---

## Task 2: Code Style Lint Fixes 📝

### Step 1: Identify Issues
```bash
# Run locally to see all issues
make lint
# or: golangci-lint run --timeout=5m
```

### Step 2: Fix Common S1009 Patterns

**Pattern:** Unnecessary nil checks on maps/slices

```go
// ❌ BEFORE (golangci-lint warning: S1009)
if m != nil && len(m) > 0 {
    for k, v := range m {
        // process
    }
}

// ✅ AFTER (no warning)
if len(m) > 0 {
    for k, v := range m {
        // process
    }
}
```

### Step 3: Find and Fix All Occurrences
```bash
# Find all S1009 warnings in code
make lint | grep S1009

# Common locations to check:
# - internal/bridge/*.go
# - internal/events/*.go
# - internal/workflows/*.go
```

### Step 4: Verify Fixes
```bash
# Run lint locally until clean
make lint

# Run tests to ensure no regressions
make test-fast
```

### Step 5: Update CI Workflow
- [ ] Edit `.github/workflows/test.yml`
- [ ] Find: `test-summary` job
- [ ] Update: Add lint back to blocking requirements
  ```yaml
  needs: [test-local, test-coverage, test-bridge, test-events, build, lint]
  
  # And add lint check to the script:
  if [[ "${{ needs.lint.result }}" == "failure" ]]; then
      echo "❌ Code style issues (golangci-lint)"
      exit 1
  fi
  ```

### Step 6: Verify in CI
- [ ] Push all changes (lint fixes + workflow update)
- [ ] Check GitHub Actions run
- [ ] Verify: `Lint` job passes ✅
- [ ] Verify: `Test Summary` passes ✅

**Definition of Done:** Lint job passes in CI for 2 consecutive runs

---

## Task 3: Integration Test Verification ✅

### Step 1: Manual Verification
```bash
# Terminal 1: Start Temporal
make run-temporal

# Terminal 2: Run integration tests
make test-with-temporal

# Should see: Test results passing
```

### Step 2: CI Verification
- [ ] All 7 jobs passing in GitHub Actions:
  - [ ] ✅ Verify Tool Versions
  - [ ] ✅ Local Tests (1.21)
  - [ ] ✅ Local Tests (1.22)
  - [ ] ✅ Tests + Coverage
  - [ ] ✅ Bridge E2E Tests
  - [ ] ✅ Event Translator Tests
  - [ ] ✅ Build Check
  - [ ] ✅ Lint
  - [ ] ✅ Integration Tests (NEW!)
  - [ ] ✅ Sprint Summary Report
  - [ ] ✅ Test Summary

### Step 3: Check Sprint Reports
- [ ] Go to GitHub Actions: https://github.com/ogormans-deptstack/orca-temporal-bridge/actions
- [ ] Find latest run
- [ ] Go to Artifacts tab
- [ ] Download `sprint-summary-report`
- [ ] Verify all tests shown as passing ✅

**Definition of Done:** All 7 jobs passing + sprint report shows green ✅

---

## Task 4: Documentation Updates 📚

### Step 1: Update README.md
- [ ] Update status section (change "Sprint 1.3 Complete" to "Sprint 1.4 Complete")
- [ ] Remove "(non-blocking)" notes from CI section
- [ ] Add note about Docker login requirement

### Step 2: Update TESTING.md (or create if missing)
```markdown
## Local Temporal Setup

### Prerequisites
- Docker installed
- Docker Hub account (or GitHub credentials)

### Steps
1. Start Temporal server: `make run-temporal`
2. Run integration tests: `make test-with-temporal`
3. View logs: `docker logs <container_id>`
```

### Step 3: Update TROUBLESHOOTING.md
Add sections:
- [ ] "Docker pull failing" → solution (docker login)
- [ ] "Lint warnings" → how to fix S1009
- [ ] "Integration tests timeout" → increase timeout or check Temporal logs

### Step 4: Update SPRINT_1_4_PLAN.md
- [ ] Mark each task as COMPLETE
- [ ] Add completion date
- [ ] Add notes/lessons learned

**Definition of Done:** All docs updated and linked from README

---

## Task 5: Branch Protection Enablement 🛡️

**Only after all CI jobs pass consistently!**

### Step 1: Verify All Jobs Pass
```bash
# Go to: https://github.com/ogormans-deptstack/orca-temporal-bridge/actions
# Confirm last 3 runs show all green ✅
```

### Step 2: Enable Branch Protection
```bash
# Option 1: Run automation script
./scripts/setup-branch-protection.sh

# Option 2: Manual GitHub UI
# Go to: Settings → Branches → Add Rule
# Branch: main
# Require all CI jobs to pass
# Require 1 review
# Dismiss stale reviews
# Require signed commits
# Include administrators: YES
```

### Step 3: Test Protection
- [ ] Create test branch
- [ ] Try to merge to main without passing CI (should block)
- [ ] Try to merge without review (should block)
- [ ] Fix CI issues, get review, merge succeeds ✅

**Definition of Done:** Branch protection active, tested, team trained

---

## Timeline & Effort

| Task | Effort | Days | Priority |
|------|--------|------|----------|
| Docker Login | 1 day | Mon | HIGH |
| Lint Fixes | 1 day | Tue | HIGH |
| Integration Verify | 0.5 day | Wed AM | MEDIUM |
| Docs Update | 0.5 day | Wed PM | MEDIUM |
| Branch Protection | 0.5 day | Thu | MEDIUM |
| **Total** | **~3 days** | **Mon-Thu** | — |

---

## Success Criteria ✅

- [ ] Docker Temporal image pulls successfully in CI
- [ ] All 7 CI jobs passing consistently (2+ consecutive runs)
- [ ] Zero lint warnings in main branch
- [ ] Sprint dashboard shows all green
- [ ] Branch protection rules enforced
- [ ] Team trained on new workflow
- [ ] Documentation complete and linked
- [ ] Ready for production deployment

---

## Rollback Plan

If docker login fails:
- [ ] Revert to `temporalapi/server:latest` from Docker Hub with credentials
- [ ] If that fails, disable integration tests and create issue for later

If lint fixes break tests:
- [ ] Revert lint changes
- [ ] Fix tests first, then reapply lint strictness

---

## Escalation

If stuck, reach out to:
- Docker issues: Check [docker-login-action](https://github.com/docker/login-action)
- Lint issues: Check [golangci-lint docs](https://golangci-lint.run/)
- Temporal issues: Check [Temporal.io troubleshooting](https://docs.temporal.io/troubleshooting)

---

**Start Date:** [TO BE FILLED]  
**Assigned To:** [TEAM MEMBER]  
**Review Date:** [TO BE FILLED]  
**Completion Date:** [TO BE FILLED]

---

**Note:** This is a continuation of Sprint 1.3. All groundwork is done. This sprint is just optimization and team readiness.
