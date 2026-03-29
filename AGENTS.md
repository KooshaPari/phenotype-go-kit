# Agent Rules

**This project is managed through AgilePlus.**

## AgilePlus Mandate

All work MUST be tracked in AgilePlus:
- Reference: `/Users/kooshapari/CodeProjects/Phenotype/repos/AgilePlus`
- CLI: `cd /Users/kooshapari/CodeProjects/Phenotype/repos/AgilePlus && agileplus <command>`

## Branch Discipline

- Feature branches in `repos/worktrees/<project>/<category>/<branch>`
- Canonical repository tracks `main` only
- Return to `main` for merge/integration checkpoints

## Work Requirements

1. **Check for AgilePlus spec before implementing** (`repos/AgilePlus/scripts/list-features.sh` or `repos/AgilePlus/kitty-specs/<slug>/`).
2. **Track delivery** with `agileplus validate --feature <slug>`, `plan`, `implement`, `queue list`, and `cycle list` as appropriate.
3. **No code without corresponding AgilePlus spec**

## UTF-8 Encoding

All markdown files must use UTF-8. Avoid smart quotes, em-dashes, and special characters where they break tooling. There is no `agileplus validate-encoding` in the current CLI; rely on editor and pre-commit.
