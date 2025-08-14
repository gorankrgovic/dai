# Commands

This page lists all available **DAI CLI** commands, their flags, and usage examples.

---

## `dai auth`

Authenticate DAI with external services (currently GitHub token only).

```bash
dai auth
```
Prompts for a GitHub Personal Access Token (PAT) and stores it in the local configuration.

---

## `dai completion`

Generate the autocompletion script for your shell.

```bash
dai completion bash
dai completion zsh
dai completion fish
```

---

## `dai config`

Configure global DAI settings stored in `~/.dai/config.yaml`.

```bash
dai config
```

---

## `dai ignore`

Create a default `.daiignore` file in the project root (uses `.gitignore` syntax).

```bash
dai ignore [flags]
```

**Flags:**
| Flag          | Description                                                   | Default         |
|---------------|---------------------------------------------------------------|-----------------|
| `--path`      | Path to create the `.daiignore` file (relative to project root) | `.daiignore`    |
| `--force`     | Overwrite existing file without prompting                     | `false`         |
| `--yes`       | Non-interactive: assume “yes” to prompts                      | `false`         |

---

## `dai init`

Initialize DAI for the current project by creating `.dai/project.yaml`.

```bash
dai init [flags]
```

**Flags:**
| Flag         | Description                                                        | Default |
|--------------|--------------------------------------------------------------------|---------|
| `-f, --force`   | Overwrite existing `.dai/project.yaml` without prompt            | `false` |
| `-v, --verbose` | Print detected git info                                          | `false` |
| `-p, --path`    | Project path (defaults to current directory)                     | `.`     |

---

## `dai triage`

Analyze a commit and open a single GitHub issue with findings.  
If `[commit]` is **not** provided, DAI will analyze the **latest commit (HEAD)** in the repository.

```bash
dai triage [commit] [flags]
```

**Examples:**
```bash
# Analyze the latest commit
dai triage

# Analyze a specific commit
dai triage 8282882

# Dry run without creating an issue
dai triage 8282882 --dry-run
```

**Flags:**
| Flag             | Description                                                        | Default                                        |
|------------------|--------------------------------------------------------------------|------------------------------------------------|
| `--ext`          | Comma-separated file extensions to analyze                         | `.js,.jsx,.ts,.tsx,.vue,.php,.py,.go`          |
| `--dry-run`      | Print the would-be GitHub issue without creating it                 | `false`                                        |
| `--model`        | Override OpenAI model from config (optional)                        | *(none)*                                       |
| `--max-file-kb`  | Max file size per analyzed file (KB)                                | `80`                                           |
| `--ignore`       | Path to ignore file (gitignore syntax), relative to project root    | `.daiignore`                                   |
| `--always-open`  | Always create a GitHub issue even when no findings                  | `false`                                        |
| `--diff-context` | Number of context lines per diff hunk                               | `3`                                            |


---

## `dai triage-local`

Analyze a single local file and append findings to `.dai/local.log`.

```bash
dai triage-local path/to/file.go [flags]
```

**Flags:**
| Flag             | Description                                          | Default             |
|------------------|------------------------------------------------------|---------------------|
| `--model`        | Override OpenAI model from config (optional)          | *(none)*            |
| `--max-file-kb`  | Max bytes per analyzed file (KB)                      | `200`               |
| `--log`          | Path to local log file (relative to project root)     | `.dai/local.log`    |
| `--format`       | Log format (`md` or `json`)                           | `md`                 |
| `--no-stdout`    | Do not print findings to stdout (log only)            | `false`             |

---

## Global Flags

| Flag                  | Description                                             |
|-----------------------|---------------------------------------------------------|
| `-h, --help`          | Show help for a command                                 |
| `-p, --parrot string` | Summon the DAI parrot (`party`, `insult`, `wise` modes) |
| `-v, --version`       | Show the DAI version                                    |

---

Next: [Troubleshooting](troubleshooting.md)
