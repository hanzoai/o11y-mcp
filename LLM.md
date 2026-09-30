# o11y-mcp

**Org:** hanzoai  ·  **Ecosystem:** hanzo  ·  **Path:** `/Users/a/work/hanzo/hanzoai/o11y-mcp`
**Origin:** https://github.com/hanzoai/o11y-mcp.git

## Discovery

This file (`CLAUDE.md`) is the canonical agent-facing readme; `LLM.md` is a symlink to it. Update either name and both stay in sync.

## Where to look first

- `README.md` — human-facing overview (if present)
- `package.json` / `Cargo.toml` / `pyproject.toml` / `go.mod` — language & deps
- `.github/workflows/` — CI surface
- `docs/` — extended docs (if present)

## Sibling repos

See the org-level `LLM.md` at `/Users/a/work/hanzo/hanzoai/LLM.md` for the full inventory of sibling repos and inter-repo dependencies.

## Routes

Every call is `/v1/o11y/…` on the o11y host `SIGNOZ_URL` names (host only, no
path), authenticated by the `O11Y-API-KEY` header o11y reads.
`internal/client/prefix_test.go` fails on any call outside `/v1/o11y/`.
