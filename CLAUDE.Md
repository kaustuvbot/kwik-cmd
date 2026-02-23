# Project Specification for OpenClaw

---

## Project Name

**kwik-cmd**

---

## Objective

Build a high-performance CLI tool written in **Go** that:

1. Tracks terminal commands automatically
2. Learns usage patterns
3. Suggests intelligent command completions
4. Supports keyword-based search
5. Maintains structured command history
6. Works natively with Bash and Zsh
7. Operates fully locally (no cloud required)

This tool should outperform shell history (`Ctrl + R`) by offering:

* Structured command understanding
* Ranking by recency + frequency
* Context awareness (directory-based)
* Flag-level intelligence

---

# Why Go Is the Best Choice

OpenClaw should use **Go**, not Python or Rust, because:

### 1️⃣ Performance

Shell hooks must execute instantly.
Go binaries are compiled and fast.

### 2️⃣ Single Binary Distribution

No runtime dependency.
Easy installation via:

* curl
* brew
* apt
* direct binary

### 3️⃣ Excellent CLI Ecosystem

Use:

* Cobra (CLI framework)
* BubbleTea (TUI)
* SQLite driver (modernc.org/sqlite or mattn/go-sqlite3)

### 4️⃣ Cross Platform

Works on:

* Linux
* macOS
* Windows

---

# Repository Requirements

OpenClaw must:

* Use existing repository
* Work strictly on `dev` branch
* Never commit to `main`
* Follow clean commit discipline
* Use conventional commits

Branch rule:

```
main → protected
dev → active development
```

---

# High-Level Architecture

```
Shell Hook (bash/zsh)
        ↓
kwik-cmd track <command>
        ↓
Parser Engine
        ↓
Metadata Extraction
        ↓
SQLite Database
        ↓
Ranking Engine
        ↓
Interactive TUI (BubbleTea)
```

---

# Core Functional Requirements

---

## 1️⃣ Command Tracking Engine

Automatically capture commands executed in terminal.

Track:

* full command
* base command
* subcommand
* flags
* arguments
* timestamp
* directory
* exit status
* frequency

Example stored record:

```
command: git commit -m "msg"
base: git
subcommand: commit
flags: [-m]
directory: /repo
success: true
frequency: 14
```

---

## 2️⃣ Suggestion Engine

When user types:

```
git
```

Show:

Recent (top 3)
Most Used (top 3)

Rank using:

Score =
(recencyWeight * timeScore) +
(frequencyWeight * usageCount) +
(directoryWeight * directoryMatch)

---

## 3️⃣ Flag Intelligence

System should understand:

```
git commit -m
```

Store:

```
flag: -m
meaning: message
keywords: [commit, message]
```

---

## 4️⃣ Keyword Search Mode

Trigger via:

```
Ctrl + K
```

User types:

```
rebase message
```

Return ranked commands matching those keywords.

Implement fuzzy search.

---

## 5️⃣ Context Awareness

Boost commands used in current directory.

Example:

Inside terraform project:
Boost terraform commands.

---

# Database Design (SQLite)

Tables:

commands

* id
* base
* subcommand
* full_command
* frequency
* last_used
* directory

flags

* id
* command_id
* flag
* meaning

keywords

* id
* command_id
* keyword

usage_stats

* id
* command_id
* success
* exit_code
* used_at

---

# Technical Requirements

Language: Go 1.22+

Libraries:

* Cobra (CLI)
* BubbleTea (TUI)
* SQLite driver
* fuzzy matcher library

Must:

* Compile to single binary
* Avoid external dependencies
* Run in < 10ms for tracking
* Run in < 50ms for suggestions

---

# CLI Commands

kwik-cmd track "<command>"

kwik-cmd suggest "<partial>"

kwik-cmd search "<keywords>"

kwik-cmd stats

kwik-cmd reset

---

# Shell Integration

For zsh:

```
preexec() {
  kwik-cmd track "$1"
}
```

For bash:

Use PROMPT_COMMAND.

---

# Development Phases

---

## Phase 1 (MVP)

* CLI structure
* SQLite setup
* Track command
* Suggest recent + most used
* Basic TUI selector

---

## Phase 2

* Flag parsing
* Keyword tagging
* Ranking engine
* Context awareness

---

## Phase 3

* Failure detection
* Pattern grouping
* Alias suggestions
* Optional AI explanation mode (future)

---

# Commit Strategy

OpenClaw should:

* Commit every logical milestone
* Keep commits small
* Follow format:

feat:
fix:
refactor:
docs:
chore:

---

# Constraints

* Work only in dev branch
* Do not modify main
* Do not introduce breaking changes
* Maintain clean modular structure
* Follow idiomatic Go

---

# Folder Structure

```
cmd/
internal/
  parser/
  db/
  tracker/
  suggester/
  tui/
pkg/
configs/
```

---

# Success Criteria

Project is successful if:

1. It tracks commands automatically
2. Suggests intelligently ranked commands
3. Works smoothly inside shell
4. Compiles to single binary
5. No performance lag
6. Codebase is modular and extensible

---

# Optional Future (Not MVP)

* AI-based command summarization
* Team history sync
* Web dashboard
* GitHub Actions plugin

---

# Final Instruction to OpenClaw

Work exclusively on `dev` branch.
Use clean architecture.
Build Phase 1 completely before moving to Phase 2.
Do not over-engineer MVP.
Keep performance first-class.
Write clear README documentation for installation and usage.

