# SuperMac — Project Instructions

macOS system management CLI tool. Modular Bash architecture.

## Architecture

```
mac (dispatcher) → lib/<category>.sh (modules)
```

- Entry point: `mac <category> <action> [args]`
- Every module sources `lib/utils.sh` for shared functions
- Every module exports `<category>_dispatch` and optionally `<category>_help`
- Global shortcuts (e.g., `mac ip` → `network:ip`) defined in `GLOBAL_SHORTCUTS` map
- `bin/mac` is a copy of the root dispatcher used for installation

## Modules

`lib/` is the canonical source. Modules: `utils`, `finder`, `display`, `wifi`, `network`, `system`, `dev`, `dock`, `audio`, `screenshot`.

## Testing

```bash
make test        # or: bash tests/test.sh
```

## Linting

```bash
make lint        # or: shellcheck lib/*.sh mac
```

## Key Constraints

- **No root-level `.sh` files** — they were deleted. Only `lib/` is canonical.
- **Do not add `set -euo pipefail` to modules** — modules rely on caller error handling from dispatcher.
- **Do not restructure `lib/utils.sh`** — it's a god file (589 lines, 8 responsibilities) but splitting it is a Phase 2 task requiring careful testing.
- **`config/config.json` is decorative** — `get_config()` is never called at runtime. Changing it has no effect.
- **`system_cleanup()` is dangerous** — single function handling 7+ destructive operations. Needs its own dedicated session with tests before any changes.
