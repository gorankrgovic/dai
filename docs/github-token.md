# GitHub Token

This guide explains how to set up a **GitHub Personal Access Token (PAT)** for DAI so it can interact with your repositories.

---

## Why do you need a GitHub token?

DAI uses the GitHub API to perform certain actions, such as:

- Creating issues directly from the CLI
- Adding default labels to issues
- Checking repository information

Without a GitHub token, these features will be disabled.

---

## Creating a GitHub Personal Access Token

1. Go to [GitHub Developer Settings](https://github.com/settings/tokens)
2. Click **"Generate new token"** → **"Generate new token (classic)"**
3. Give your token a **descriptive name** (e.g., `DAI CLI Tool`)
4. Set an **expiration date** (recommended: 90 days for security)
5. Select the following scopes:
   - `repo` — Full control of private repositories (required for issue creation)
6. Click **Generate token** and **copy it** — you will not be able to see it again!

---

## Adding the token to DAI

You can add the token by running:

```bash
dai github-token
```

Or by editing your configuration file manually:

```yaml
# ~/.dai/config.yaml
github_token: ghp_yourtokenhere
```

---

## Security tips

- **Never share** your token with anyone
- Treat it like a password
- Revoke it immediately if you suspect it’s been compromised:
  [Revoke token](https://github.com/settings/tokens)

---

Next: [Commands](commands.md)
