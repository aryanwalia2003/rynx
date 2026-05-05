# 🧠 Rynx: Core Philosophy

* [ ] **Zero config to start** (works after install)
* [ ] **Works with existing Git + Jira setup** (no migration)
* [ ] **Sensible defaults** everywhere
* [ ] **Interactive + non-interactive** modes
* [ ] **Fast** (<200ms for most commands)
* [ ] **Works offline** where possible
* [ ] **Idempotent commands** (safe to re-run)
* [ ] **Clear undo / rollback** for actions

---

# 🚀 1. Onboarding & Setup

* [ ] `dev init` auto-detects:
  * git repo
  * remote (GitHub/GitLab)
  * Jira project
* [ ] **OAuth login** (no API tokens copy-paste hell)
* [ ] **Config auto-generated** (`.devrc`)
* [ ] **Team presets** (load org defaults)
* [ ] **Self-check command**: `dev doctor`
* [ ] **First-run interactive tutorial** (skippable)
* [ ] **Works with mono-repos**

---

# 🔗 2. Git + Jira Linking

* [ ] **Auto-detect Jira ticket** from:
  * branch name
  * commit message
* [ ] **Smart branch naming**: `dev start ZFW-123` → `feat/zfw-123-add-trip-logs`
* [ ] **Auto-assign Jira ticket**
* [ ] **Auto-transition status**:
  * start → In Progress
  * PR → In Review
  * merge → Done
* [ ] **Backfill linking** for old branches
* [ ] **Multi-ticket support**

---

# ⚡ 3. Task Execution Flow

### Start work
* [ ] `dev start <ticket>`: creates branch, pulls latest main, assigns ticket, shows summary

### During work
* [ ] `dev status`: shows current ticket, branch, PR status, CI status
* [ ] **Smart commit helper**: `dev commit` (suggests message, auto-tags ticket)

### Finish work
* [ ] `dev ship`: push branch, create PR, update Jira, assign reviewers

---

# 🧠 4. Intelligence Layer

* [ ] **PR title + description** auto-generated
* [ ] **Diff summarization**
* [ ] **Change categorization**: feat / fix / refactor / chore
* [ ] **Risk scoring**: large diff, critical files touched
* [ ] **Reviewer suggestion** (git blame based)
* [ ] **Detect missing tests**
* [ ] **Detect debug logs / console prints**
* [ ] **Detect secrets accidentally added**
* [ ] **Suggest commit squashing**

---

# 🔍 5. Review & PR Experience

* [ ] `dev pr view` (terminal-friendly)
* [ ] **Inline comments** from CLI
* [ ] **Approve / request changes**
* [ ] **Show CI status + logs**
* [ ] **Show diff summary** (not raw diff spam)
* [ ] **Rebase / merge** from CLI
* [ ] **Auto-resolve simple conflicts**

---

# 🧪 6. CI/CD Awareness

* [ ] **Detect failing tests** before PR
* [ ] **Run tests locally** before ship
* [ ] **Show CI results** in CLI
* [ ] **Retry failed jobs**
* [ ] **Block PR** if critical checks fail
* [ ] **Deployment status visibility**

---

# ⚙️ 7. Git Hygiene

* [ ] **Detect branch behind main**
* [ ] **Auto-rebase option**
* [ ] **Warn about messy commit history**
* [ ] **Squash helper**
* [ ] **Detect large binary files**
* [ ] **Detect accidental commits** (node_modules etc.)
* [ ] **Clean branch after merge**

---

# 🧬 8. Terminal UX

* [ ] **Beautiful TUI mode** (optional)
* [ ] **Keyboard-first navigation**
* [ ] **Fuzzy search** (like fzf)
* [ ] **Color-coded statuses**
* [ ] **Minimal output** (no noise)
* [ ] **Verbose mode** when needed
* [ ] **Smart suggestions**: “Did you mean dev ship?”
* [ ] **Copy-paste friendly output**

---

# 🧩 9. Customization & Extensibility

* [ ] **Config file** (`.devrc`)
* [ ] **Plugin system**
* [ ] **Custom workflows**: `dev flow deploy`
* [ ] **Hooks**: pre-commit, pre-PR
* [ ] **Aliases**: `dev s` → `dev start`
* [ ] **Team-wide shared configs**

---

# 🔐 10. Security & Safety

* [ ] **Secret masking**
* [ ] **Permission checks** before actions
* [ ] **Dry-run mode**: `dev ship --dry-run`
* [ ] **Confirmation** for destructive ops
* [ ] **Audit log of actions**

---

# 📊 11. Insights & Analytics

* [ ] **PR turnaround time**
* [ ] **Review delays**
* [ ] **Dev productivity stats**
* [ ] **Bottleneck detection**
* [ ] **“stuck ticket” alerts**

---

# 🌐 12. Multi-Platform & Ecosystem

* [ ] **GitHub + GitLab support**
* [ ] **Jira + Linear support** (future-proof)
* [ ] **Works on Mac/Linux/Windows**
* [ ] **CI integration** (GitHub Actions etc.)
* [ ] **Slack/Discord notifications**

---

# 🧨 13. Small Details

* [ ] **Auto-open PR in browser** (optional)
* [ ] **Clipboard integration**
* [ ] **Remember last used ticket**
* [ ] **Smart defaults** per repo
* [ ] **Instant feedback** (no waiting)
* [ ] **Graceful error messages** (no stack dumps)
* [ ] **“Fix it for me” suggestions**
