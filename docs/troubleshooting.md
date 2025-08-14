# Troubleshooting

This page lists common issues you might encounter when using **DAI** and how to fix them.

---

## Command not found: `dai`

**Cause:** The DAI binary is not in your system `PATH`.

**Solution:**
```bash
sudo mv dai /usr/local/bin
```
Then verify:
```bash
dai --version
```

---

## `openai_key` missing or invalid

**Cause:** Your OpenAI API key is not set or is incorrect.

**Solution:**
1. Run:
```bash
dai config
```
2. Enter a valid OpenAI API key (starts with `sk-`).
3. If the issue persists, check your `~/.dai/config.yaml` file.

---

## GitHub issue creation not working

**Cause:** Missing or invalid GitHub Personal Access Token (PAT).

**Solution:**
1. Generate a new token following [GitHub Token setup](github-token.md).
2. Add it to DAI:
```bash
dai github-token
```

---

## `permission denied` when running DAI

**Cause:** The binary does not have execution permissions.

**Solution:**
```bash
chmod +x dai
```

---

## Model not supported

**Cause:** You specified a model that the current OpenAI API does not support.

**Solution:**
- Update your `model` in `~/.dai/config.yaml` to a supported one (e.g., `gpt-4o-mini`).
- Check OpenAI's [model list](https://platform.openai.com/docs/models).

---

## Still having issues?

- Update DAI to the latest version:
```bash
git pull origin main && go build -o dai && sudo mv dai /usr/local/bin
```
- Open a GitHub Issue with details:
  [https://github.com/gorankrgovic/dai/issues](https://github.com/gorankrgovic/dai/issues)

---

Next: [Contributing](contributing.md)
