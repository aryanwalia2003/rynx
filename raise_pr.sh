#!/bin/bash

# Ensure gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed. Please install it first: https://cli.github.com/"
    exit 1
fi

echo "🚀 Let's raise a good PR!"
echo "--------------------------"

# Determine base branch
base_branch="$1"
if [ -z "$base_branch" ]; then
    read -p "Base branch to point this PR to (default: main): " base_branch
    base_branch=${base_branch:-main}
fi
echo "Pointed to base branch: $base_branch"
echo "--------------------------"

# Function to read multi-line input
read_multiline() {
    local prompt="$1"
    local result=""
    echo -e "\n📝 $prompt"
    echo "(Type your content. Press Enter on an empty line to finish this section)"
    while read -r line; do
        if [ -z "$line" ]; then
            break
        fi
        # Append line to result
        if [ -z "$result" ]; then
            result="$line"
        else
            result="$result\n$line"
        fi
    done
    echo -e "$result"
}

# 1. Title
echo ""
read -p "📌 PR Title (Short actionable title): " title

# 2. Summary
summary=$(read_multiline "Summary: What does this PR do? (2-5 bullet points max)")

# 3. Why is this change needed?
why_needed=$(read_multiline "Why is this change needed? (Describe the actual problem: Bug, Limitation, Scaling issue, etc.)")

# 4. What changed?
what_changed=$(read_multiline "What changed? (Describe implementation details, architecture notes)")

# 5. Reviewer Guide
reviewer_guide=$(read_multiline "Reviewer Guide: Tell reviewers where to focus (e.g. main review areas, low priority)")

# 6. How to test
how_to_test=$(read_multiline "How to test: Provide exact reproducible steps (e.g. Local \`make run\`)")

# Construct the body
body="## Summary
$summary

---

## Why is this change needed?
$why_needed

---

## What changed?
$what_changed

---

## Reviewer Guide
$reviewer_guide

---

## How to test
$how_to_test
"

# Save a backup of the body just in case the PR creation fails
backup_file="/tmp/pr_body_backup_$$.md"
echo -e "$body" > "$backup_file"

echo -e "\n========================================="
echo "Raising PR with title: $title"
echo "Target branch: $base_branch"
echo "========================================="

# Execute the gh PR create command
gh pr create --title "$title" --body "$body" --base "$base_branch"

if [ $? -eq 0 ]; then
    echo "✅ PR created successfully!"
    rm -f "$backup_file"
else
    echo "❌ Failed to create PR."
    echo "Your PR body has been saved to: $backup_file"
fi
