# dspo

**dspo** (Desktop Profile Orchestrator) is a local-first tool for capturing and replaying Linux desktop configuration intent across machines.

It focuses on **safe, additive reproduction** of a workstation environment using existing system tools (package managers, desktop settings, dotfiles), rather than attempting full system determinism.

---

## Goals

* Reduce the time required to set up a new Linux workstation
* Make desktop configuration **repeatable, versioned, and reviewable**
* Preserve existing imperative workflows (no “rewrite your OS” requirement)

---

## Non-Goals

* Full system determinism or binary reproducibility
* Replacing Nix, NixOS, or home-manager
* Managing servers, clusters, or cloud infrastructure
* Automatic destructive cleanup or reconciliation

---

## Scope (v1)

* Linux desktop environments
* GNOME only (initially)
* Ubuntu LTS and Fedora
* Safe, additive application of configuration

---

## Status

**Early development** — APIs, formats, and behavior are unstable.

This project is currently developed for personal use.
Documentation and public release may come later.

---

## License

TBD
