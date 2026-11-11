# Temporal Setup for orca-temporal-bridge

## Summary: What Works Locally

✅ **WORKS:** `temporal server start-dev` (Official CLI)
- **Pros:** Fast, no Docker, no config needed, SQLite built-in
- **Cons:** Only for local development
- **Status:** ✅ Tested and working

❓ **PARTIALLY:** Docker Compose with `temporalio/auto-setup`
- **Issue:** Requires Postgres/MySQL/Cassandra (no SQLite support in container)
- **Workaround:** Use PostgreSQL container + auto-setup
- **Status:** Complex, not recommended for local

❌ **DOES NOT WORK:** `ghcr.io/temporalapi/server`
- **Issue:** Requires authentication (GHCR_TOKEN)
- **Status:** Blocked in CI/CD

---

## Recommended Solution for orca-temporal-bridge

### For Local Development (RECOMMENDED)

**Use the Temporal CLI:**
```bash
# Install
brew install temporal-cli

# Run dev server (no Docker needed)
temporal server start-dev

# In another terminal, tests will automatically connect:
go test ./...

# View Temporal UI
open http://localhost:8233/ui
```

**Pros:**
- ✅ Zero configuration
- ✅ Built-in SQLite
- ✅ Fast startup (~5s)
- ✅ Same as production Temporal setup
- ✅ Works with our Go client (localhost:7233)

### For CI/CD (Docker)

**Option A: Use Postgres (Recommended)**
```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: temporal
      POSTGRES_PASSWORD: temporal
      POSTGRES_DB: temporal

  temporal:
    image: temporalio/auto-setup:latest
    environment:
      DB: postgres12
      POSTGRES_SEEDS: postgres
    depends_on:
      postgres:
        condition: service_healthy
```

**Option B: Pre-install Temporal CLI in Docker**
```dockerfile
FROM golang:1.22
RUN brew install temporal-cli
CMD ["temporal server start-dev"]
```

---

## Action Items for Sprint 1.4

### For PR #2 (docker-login)

**Recommended Change:**
```yaml
# .github/workflows/test.yml - Integration Tests job

# Instead of:
# services:
#   temporal:
#     image: ghcr.io/temporalapi/server:latest

# Use one of:
# 1. Postgres approach (above) - safer for CI
# 2. Remove Docker entirely - use `brew install temporal-cli` in CI runner
```

### For Repository Documentation

**Add to SETUP.md:**
```markdown
## Running Temporal Locally

### Quick Start (Recommended)
```bash
brew install temporal-cli
temporal server start-dev
```

### With Docker
```bash
docker-compose up
```
```

---

## Test Results

✅ Local tests pass with `temporal server start-dev`:
```
github.com/ogormans-deptstack/orca-temporal-bridge/internal/bridge     OK (0.444s)
github.com/ogormans-deptstack/orca-temporal-bridge/internal/events     OK (0.327s)
```

✅ CLI connects successfully:
```
$ temporal workflow list
(no workflows, as expected)
```

✅ Web UI is accessible:
```
http://localhost:8233/ui
```

---

## Files Updated

- ✅ `docker-compose.yml` - Documented recommended approaches with comments
