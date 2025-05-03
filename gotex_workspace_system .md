# gotex Module System

## File Types
- `gotex.work` – Optional workspace config (multi-document setup)
- `gotex.doc` – Defines a document (metadata, dependencies, entry file)
- `gotex.mod` – Defines a reusable module (layout, logic, themes)
- `gotex.sum` – Lockfile for resolved versions and checksums

## Implicit Behavior
- A **workspace** always exists:
  - If `gotex.work` is missing, a workspace named `default` is created implicitly
  - The root folder is its primary location, but not its name
- The **root folder** can be either a **document** or a **module**, but not both
- A **document** always exists:
  - If `gotex.doc` is missing, a document named `default` is created implicitly
  - The `.tex` file passed to `gotex <file>.tex` becomes the document’s **entry file**
  - The document is automatically registered in the implicit workspace
- A **module** is any folder that:
  - Contains a `gotex.mod`, or
  - Contains one or more `.tex` files
  - Modules are resolved only when referenced or explicitly listed in `gotex.work`

  
design ideas:
- Logical Import Names – Use registry-based names instead of VCS URLs.
- Clean Import Paths – Remove version numbers from import paths.
- Independent Package Versioning – Allow sub-packages to version separately within a module.
- Flexible Version Resolution – Prefer highest compatible versions with conflict resolution.
- Transitive Constraints Support – Let dependencies specify their own version bounds.
- Pluggable Registry System – Support custom, decentralized registries.
- Compact Checksum Management – Reduce or segment go.sum-like data.
- Explicit Module Boundaries – Define modules via config, not directory inference.
- Optional & Scoped Dependencies – Support test/dev/runtime scopes.
- Better Fork Handling – Simplify switching and maintaining forks in builds.



---
# Gotex Module & Build System — Summary

## Primitives

- Primitives are built into the Gotex engine.
- Always available without import (`\section`, `\chapter`, `\matrix`, `\grid`, etc.).
- No modules or configuration needed to use them.
- Offline by default.

## Modules

- A module is any folder with `gotex.mod` or imported via `\importmodule{...}`.
- All imported folders are treated as modules (like Go packages).
- Implicit `gotex.mod` assumed if not present.
- Modules can contain `.tex` files, macros, styles, and assets.
- Versioned and reusable when declared in `gotex.doc`.

## Config Files

- `gotex.work` (optional): multi-doc workspace config, registry overrides.
- `gotex.doc` (implicit if missing): defines entry file, Gotex version, and module dependencies.
- `gotex.sum`: created only when external modules are used; locks versions and checksums.

## Module Resolution

- Resolution order:
  1. `replace` entries in config
  2. Local folders (internal modules)
  3. Registry fetch (external modules)
- Import syntax is clean: `\importmodule{matrix}`, no version or source in path.
- External modules are zipped (`.gtxmod.zip`), downloaded and cached.
- External usage auto-generates `gotex.doc` and `gotex.sum`.

## Resource Handling

- Resources (images, fonts, data) are **not modules**.
- Referenced directly via relative paths: `\image{assets/logo.png}`.
- Resolved from filesystem, not module system.
- No versioning or config needed unless bundled in a module explicitly.

## Philosophy

- All behavior is consistent, whether implicit or explicit.
- Engine-first, offline-first, zero-config start.
- Grows into full modular system as needed.
- Module system supports versioning, reuse, and external distribution cleanly.

---

# Gotex Build Plan Spec (Local Module Only, Initial Phase)

## 1. Target

- A `.tex` file passed to `gotex <file>.tex`
- Serves as entry point for rendering
- Must resolve to or create a `gotex.doc`

## 2. Build Context

- If `gotex.work` is found: load workspace definition
- If missing: create implicit workspace with name `"default"`
- If `gotex.doc` is found: load and validate
- If missing: create implicit document
  - `name = "default"`
  - `entry = <target>`
  - `requires = []`
- All relative paths resolved from the location of `gotex.doc`

## 3. Scan Target

- Confirm `entry` file exists
- Load only enough to check syntax and detect imports

## 4. Parse ImportsOnly (Target)

- Tokenize `target.tex`
- Collect all `\input{}`, `\include{}`, `\usemodule{}` statements
- No macro expansion, no execution

## 5. Validate Imports

- For each `\usemodule{...}`:
  - Check `requires[]` list in `gotex.doc`
  - Match symbolic name (e.g., `"layout.invoice"`) against `provides[]` in `gotex.mod`
  - If unresolved → error

## 6. Scan Import Folders

- For each path in `requires[]`:
  - Must exist and be a folder
  - Must contain either:
    - A `gotex.mod`, or
    - At least one `.tex` file

## 7. Parse ImportsOnly (Imports)

- Tokenize `.tex` files in each required module (ImportsOnly mode)
- Collect all additional `\usemodule{}`, `\input{}`, `\include{}` statements
- Resolve relative to module root

## 8. Repeat (5–7) Until Completion

- Resolve all imports recursively
- Track visited paths/modules
- Detect cycles in document/module graph
- Report cycle errors with a clear path trace

## 9. DAG Construction

- Build dependency DAG:
  - Nodes = documents/modules/files
  - Edges = includes, module uses
- Sort topologically
- Mark target document as root node
- Pass ordered list to full parser/layout engine


---

# Gotex Project Package Structure:

gotex/
├── cmd/                # CLI entrypoint: parses flags, invokes build
├── config/             # Parses gotex.doc, gotex.mod, gotex.work (TOML to structs)
├── module/             # Resolves module paths, metadata, and provides[] symbols
├── document/           # Handles document files, entry.tex, structure, \inputs
├── build/              # Coordinates full build: context, scan, ImportsOnly, DAG
├── engine/             # Full rendering pipeline (IR generation, layout, emission)
├── token/              # Pos, Position, FileSet, Token enums (already in place)
├── scanner/            # Tokenizer: input stream to tokens
├── parser/             # Builds syntax tree or IR from token stream

# With Go compiler references

gotex/
├── cmd/                # → `cmd/go/`, `cmd/compile/` – CLI entry, flag handling
├── config/             # → `cmd/go/internal/modfile`, `cmd/go/internal/work` – config file parsing
├── module/             # → `cmd/go/internal/load`, `go/build`, `cmd/go/internal/modload` – resolves modules, deps
├── document/           # → `cmd/go/internal/load`, `cmd/go/internal/work` – input file orchestration
├── build/              # → `cmd/go/internal/work` – BuildPlan, DAG, dependency graph, scheduler
├── engine/             # → `cmd/compile/internal/gc`, `go/types`, `go/ssa` – actual processing/rendering
├── token/              # → `go/token` – Pos, FileSet, token types
├── scanner/            # → `go/scanner`, `cmd/compile/internal/syntax/scanner.go`
├── parser/             # → `go/parser`, `cmd/compile/internal/syntax/parse.go`

---

# Order of development

token/ – foundational definitions (positions, tokens) (done)

scanner/ – depends only on token

parser/ – depends on scanner and token

config/ – config parsing; logically independent, needs basic types

module/ – needs config and token, optionally parser

document/ – needs config, token, scanner, maybe parser

build/ – orchestrates above, depends on config, module, document

engine/ – final layout and rendering, depends on parser, document, etc.

cmd/ – top-level CLI, integrates build, config, engine
