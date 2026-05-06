package prcmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"rynx/internal/config"
	"rynx/shared/configutils"
	"rynx/shared/errors"
	"rynx/shared/promptutils"
)

const (
	ansi_bold   = "\033[1m"
	ansi_reset  = "\033[0m"
	ansi_cyan   = "\033[36m"
	ansi_yellow = "\033[33m"
	ansi_green  = "\033[32m"
	ansi_dim    = "\033[2m"
)

func (p *PR_struct) PR_run_method(args []string) error {
	// ── ensure github token ─────────────────────────────────────────────────
	path, err := configutils.Config_global_path_util()
	if err != nil {
		return err
	}
	cfg, _ := configutils.Config_load_util(path)
	if cfg == nil {
		cfg = &config.Config_struct{}
	}

	if cfg.Git.GitHubToken == "" {
		fmt.Printf("\n%sGitHub token not found.%s\n", ansi_yellow, ansi_reset)
		fmt.Println("Create a Personal Access Token with 'repo' scope at: https://github.com/settings/tokens")
		token, err := promptutils.Prompt_string_util("Enter GitHub Token", "")
		if err != nil {
			return err
		}
		if strings.TrimSpace(token) == "" {
			return errors.Error_wrap_util("GitHub token is required to raise a PR", nil)
		}
		cfg.Git.GitHubToken = strings.TrimSpace(token)
		if err := configutils.Config_save_util(path, cfg); err != nil {
			return errors.Error_wrap_util("Failed to save config", err)
		}
		fmt.Printf("%sToken saved globally!%s\n", ansi_green, ansi_reset)
	}

	// ── resolve base branch ─────────────────────────────────────────────────
	baseBranch := ""
	if len(args) > 0 {
		baseBranch = args[0]
	}
	if baseBranch == "" {
		var err error
		baseBranch, err = promptutils.Prompt_string_util("Base branch to merge into", "main")
		if err != nil {
			return err
		}
	}

	printHeader("Raising a PR → " + baseBranch)

	// ── title ────────────────────────────────────────────────────────────────
	title, err := promptutils.Prompt_string_util(
		ansi_bold+"📌  PR Title"+ansi_reset+" (short, actionable — e.g. 'Add retry support for webhook delivery')",
		"",
	)
	if err != nil {
		return err
	}
	if strings.TrimSpace(title) == "" {
		return errors.Error_wrap_util("PR title cannot be empty", nil)
	}

	// ── summary ──────────────────────────────────────────────────────────────
	summary, err := promptutils.Prompt_multiline_util(
		ansi_bold + "📋  Summary" + ansi_reset + " — What does this PR do?" + ansi_dim + " (2–5 bullet points)" + ansi_reset,
	)
	if err != nil {
		return err
	}

	// ── why needed ───────────────────────────────────────────────────────────
	whyNeeded, err := promptutils.Prompt_multiline_util(
		ansi_bold + "🤔  Why is this change needed?" + ansi_reset + ansi_dim + " (bug, limitation, product requirement…)" + ansi_reset,
	)
	if err != nil {
		return err
	}

	// ── what changed ─────────────────────────────────────────────────────────
	whatChanged, err := promptutils.Prompt_multiline_util(
		ansi_bold + "🔧  What changed?" + ansi_reset + ansi_dim + " (implementation details, architecture notes)" + ansi_reset,
	)
	if err != nil {
		return err
	}

	// ── reviewer guide ───────────────────────────────────────────────────────
	reviewerGuide, err := promptutils.Prompt_multiline_util(
		ansi_bold + "👀  Reviewer Guide" + ansi_reset + ansi_dim + " (key files, focus areas, low-priority items)" + ansi_reset,
	)
	if err != nil {
		return err
	}

	// ── how to test ──────────────────────────────────────────────────────────
	howToTest, err := promptutils.Prompt_multiline_util(
		ansi_bold + "🧪  How to test?" + ansi_reset + ansi_dim + " (exact reproducible steps)" + ansi_reset,
	)
	if err != nil {
		return err
	}

	// ── build body ──────────────────────────────────────────────────────────
	body := buildBody(summary, whyNeeded, whatChanged, reviewerGuide, howToTest)

	// ── confirm ─────────────────────────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════%s\n", ansi_cyan, ansi_reset)
	fmt.Printf("%s  Title:%s %s\n", ansi_bold, ansi_reset, title)
	fmt.Printf("%s  Base: %s %s\n", ansi_bold, ansi_reset, baseBranch)
	fmt.Printf("%s══════════════════════════════════════════%s\n\n", ansi_cyan, ansi_reset)

	confirm, err := promptutils.Prompt_string_util("Raise PR? [Y/n]", "Y")
	if err != nil {
		return err
	}
	if strings.ToLower(strings.TrimSpace(confirm)) == "n" {
		fmt.Println("Aborted.")
		return nil
	}

	// ── native github api call ──────────────────────────────────────────────
	// 1. Get current branch
	branchOut, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return errors.Error_wrap_util("Failed to get current git branch", err)
	}
	currentBranch := strings.TrimSpace(string(branchOut))
	if currentBranch == "" {
		return errors.Error_wrap_util("Not currently on any branch", nil)
	}

	// 2. Get remote origin URL
	remoteOut, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return errors.Error_wrap_util("Failed to get remote origin URL", err)
	}
	remoteURL := strings.TrimSpace(string(remoteOut))

	// Parse owner/repo from remote URL (supports HTTPS and SSH)
	// Example SSH: git@github.com:aryanwalia2003/rynx.git
	// Example HTTPS: https://github.com/aryanwalia2003/rynx.git
	repoPath := ""
	if strings.HasPrefix(remoteURL, "git@github.com:") {
		repoPath = strings.TrimPrefix(remoteURL, "git@github.com:")
	} else if strings.HasPrefix(remoteURL, "https://github.com/") {
		repoPath = strings.TrimPrefix(remoteURL, "https://github.com/")
	} else {
		return errors.Error_wrap_util("Unsupported git remote format. Must be GitHub.", nil)
	}
	repoPath = strings.TrimSuffix(repoPath, ".git")

	// 3. Push branch to remote
	fmt.Printf("\n%sPushing %s to origin...%s\n", ansi_dim, currentBranch, ansi_reset)
	pushCmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return errors.Error_wrap_util("Failed to push branch to remote", err)
	}

	// 4. Create PR via GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repoPath)
	payload := map[string]string{
		"title": title,
		"body":  body,
		"head":  currentBranch,
		"base":  baseBranch,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errors.Error_wrap_util("Failed to create request", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Git.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Error_wrap_util("GitHub API request failed", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		json.Unmarshal(respBody, &errResp)
		errMsg := string(respBody)
		if msg, ok := errResp["message"].(string); ok {
			errMsg = msg
		}
		return errors.Error_wrap_util("Failed to create PR: "+errMsg, nil)
	}

	var successResp struct {
		HTMLURL string `json:"html_url"`
	}
	json.Unmarshal(respBody, &successResp)

	fmt.Printf("\n%s✅  PR raised successfully!%s\n", ansi_green, ansi_reset)
	fmt.Printf("🔗  %s\n", successResp.HTMLURL)
	return nil
}

// ── helpers ─────────────────────────────────────────────────────────────────

func printHeader(msg string) {
	fmt.Printf("\n%s%s🚀  %s%s\n", ansi_bold, ansi_cyan, msg, ansi_reset)
	fmt.Println(strings.Repeat("─", 50))
}

func buildBody(summary, whyNeeded, whatChanged, reviewerGuide, howToTest string) string {
	var b strings.Builder

	section := func(heading, content string) {
		b.WriteString("## ")
		b.WriteString(heading)
		b.WriteString("\n\n")
		if strings.TrimSpace(content) == "" {
			b.WriteString("_N/A_")
		} else {
			b.WriteString(content)
		}
		b.WriteString("\n\n---\n\n")
	}

	section("Summary", summary)
	section("Why is this change needed?", whyNeeded)
	section("What changed?", whatChanged)
	section("Reviewer Guide", reviewerGuide)
	section("How to test", howToTest)

	return b.String()
}
