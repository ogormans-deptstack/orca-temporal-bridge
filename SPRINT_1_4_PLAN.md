# Sprint 1.4 Plan: Docker Login & Lint Strictness

**Status:** Planned (after Sprint 1.3 completes)  
**Priority:** Medium  
**Estimated effort:** 2-3 days  

---

## Overview

Sprint 1.3 successfully established CI/CD with version management and branch protection, but two optional jobs are still skipped:
- **Integration Tests** (Temporal Docker pull failing)
- **Lint Job** (Code style warnings blocking)

Sprint 1.4 will bring these fully online before team scaling.

---

## Task 1: Docker Hub Login for Temporal Integration Tests

**Current issue:** `ghcr.io/temporalapi/server:latest` pull fails with 403 access denied (public rate limits).

### Solution: GitHub Actions Docker Login

1. **Create GitHub Personal Access Token (PAT)**
   ```bash
   # On GitHub: Settings → Developer settings → Personal access tokens → Tokens (classic)
   # Required scopes: read:packages
   # Save token as DOCKERHUB_TOKEN in your notes (you'll only see it once)
   ```

2. **Add GitHub Secrets**
   ```bash
   # In repo: Settings → Secrets and variables → Actions
   # Add two secrets:
   # - DOCKERHUB_USERNAME: your GitHub username (or Docker Hub username if using Docker Hub)
   # - DOCKERHUB_TOKEN: your PAT token
   ```

3. **Update CI Workflow**
   ```yaml
   # In .github/workflows/test.yml, add to test-integration job:
   
   - name: Log in to Docker
     uses: docker/login-action@v2
     with:
       registry: ghcr.io
       username: ${{ github.actor }}
       password: ${{ secrets.GITHUB_TOKEN }}
   ```

4. **Alternative: Use Docker Hub Credentials**
   ```yaml
   - name: Log in to Docker Hub
     uses: docker/login-action@v2
     with:
       username: ${{ secrets.DOCKERHUB_USERNAME }}
       password: ${{ secrets.DOCKERHUB_TOKEN }}
   
   # Then use: image: temporalapi/server:latest (back to original)
   ```

### Testing Locally
```bash
# Simulate the integration test locally
docker pull ghcr.io/temporalapi/server:latest
docker-compose up -d  # Or: docker run temporalapi/server
make test-with-temporal
```

### Definition of Done
- [ ] Secrets configured in GitHub repo
- [ ] Docker login step added to test-integration job
- [ ] CI run succeeds with all 7 jobs passing
- [ ] No more "pull access denied" errors in logs

---

## Task 2: Code Style Lint - Enable Strict Mode

**Current issue:** golangci-lint finds S1009 issues (unnecessary nil checks on maps/slices).

### Solution: Fix Code + Enable Linting

1. **Review Lint Warnings**
   ```bash
   # Run locally to see all issues
   make lint  # or: golangci-lint run --timeout=5m
   ```

2. **Fix Common S1009 Issues**
   
   Example fix:
   ```go
   // BEFORE (lint warning)
   if m != nil && len(m) > 0 {
       // use map
   }
   
   // AFTER (no nil check needed)
   if len(m) > 0 {
       // use map - Go's len() returns 0 for nil maps
   }
   ```

3. **Make Lint Blocking Again**
   ```yaml
   # In test-summary job, add back lint validation:
   if [[ "${{ needs.lint.result }}" == "failure" ]]; then
       echo "❌ Code style issues (golangci-lint)"
       exit 1
   fi
   ```

4. **Update Dependencies**
   ```bash
   # Ensure latest golangci-lint in CI
   # In .github/workflows/test.yml:
   uses: golangci/golangci-lint-action@v4  # Consider upgrading
   ```

### Testing Locally
```bash
make lint           # See all issues
make fmt            # Auto-format code
make vet            # Additional checks
```

### Definition of Done
- [ ] All S1009 warnings fixed in code
- [ ] Local `make lint` passes cleanly
- [ ] Lint job added back to test-summary's blocking requirements
- [ ] CI passes all 7 jobs with lint ✓

---

## Task 3: Integration Tests - Temporal Full Run

**Once docker login works**, verify full integration:

1. **Manual test run**
   ```bash
   make run-temporal        # Start Temporal server
   make test-with-temporal  # Run integration tests
   ```

2. **CI verification**
   ```bash
   # Push to develop branch, check GitHub Actions
   # All 7 jobs should pass:
   # ✓ verify-versions
   # ✓ test-local (1.21)
   # ✓ test-local (1.22)
   # ✓ test-coverage
   # ✓ test-bridge
   # ✓ test-events
   # ✓ test-integration    <-- Now working
   # ✓ build
   # ✓ lint                <-- Now blocking
   ```

### Definition of Done
- [ ] test-integration job passes in CI
- [ ] Sprint summary report shows all green ✅
- [ ] No warnings in GitHub Actions logs

---

## Task 4: Update Documentation

Once lint and integration tests pass, update:

1. **README.md**
   - Remove "(non-blocking)" note from CI section
   - Confirm all 7 jobs documented

2. **TESTING.md**
   - Add docker login setup instructions
   - Document local Temporal server startup

3. **TROUBLESHOOTING.md**
   - Add "Docker pull failing" → solution section
   - Add "Lint warnings" → fixes section

---

## Rollout & Branch Protection

Once Sprint 1.4 complete:

```bash
# Enable full branch protection
./scripts/setup-branch-protection.sh

# Now all merges to main require:
# ✅ All 7 CI jobs passing
# ✅ 1 code review
# ✅ Resolved conversations
# ✅ Signed commits
```

---

## Risk & Mitigation

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Docker token expires | Integration tests fail | Store in GitHub Secrets, rotate quarterly |
| Lint strictness breaks workflow | Developer friction | Start with warnings-only, gradually strict |
| Temporal pulls slow | CI timeout | Cache Docker images, use local registry |

---

## Success Metrics

- ✅ All 7 GitHub Actions jobs pass consistently
- ✅ Integration tests run in <1 minute
- ✅ Zero lint warnings in CI
- ✅ Sprint summary reports all green
- ✅ Ready to share with team colleagues

---

## Next Steps

1. **Assign to team member** who has Docker Hub/GitHub access
2. **Create PR** with docker login changes
3. **Test locally first** before pushing to CI
4. **Monitor first run** for any issues
5. **Celebrate** 🎉 when all 7 jobs pass!

---

## Related Files

- `.github/workflows/test.yml` — Main CI workflow
- `Makefile` — Local test targets
- `docker-compose.yml` — Temporal server config
- `.github/dependabot.yml` — For automated dependency updates (future)

---

## Timeline Estimate

- **Task 1 (Docker Login):** 1 day (setup secrets, test, verify)
- **Task 2 (Lint Fixes):** 1 day (fix code, verify locally)
- **Task 3 (Integration):** 0.5 day (run full tests, document)
- **Task 4 (Docs):** 0.5 day (update guides)

**Total:** ~3 days of focused work  
**Blocking:** Nothing (can happen in parallel with other Sprint 1.4 work)
