# Configuration

This guide explains how to configure **DAI** for the first time and how to update your settings later.

---

## Local configuration

DAI stores its configuration locally in:

```text
~/.dai/config.yaml
```

This file is automatically created when you run:

```bash
dai config
```

---

## Configuration fields

| Field          | Description                                                                 | Example                               |
|----------------|-----------------------------------------------------------------------------|---------------------------------------|
| `openai_key`   | Your OpenAI API key (required for AI features)                              | `sk-1234567890abcdef`                 |
| `model`        | Preferred AI model                                                          | `gpt-4o-mini`                         |
| `provider`     | (Optional) API provider name                                                | `openai`                              |
| `github_token` | (Optional) GitHub Personal Access Token for GitHub integration              | `ghp_1234567890abcdef`                 |

---

## Editing the config manually

You can edit the config file directly using your preferred text editor:

```bash
nano ~/.dai/config.yaml
```

Make sure to keep correct YAML formatting.

---

## Resetting configuration

If you want to start fresh:

```bash
rm ~/.dai/config.yaml
dai config
```

This will recreate the config file and start the setup wizard again.

---

Next: [GitHub Token](github-token.md)
