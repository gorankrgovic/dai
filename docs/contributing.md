# Contributing

We welcome contributions to **DAI**!  
Whether it's fixing a bug, improving documentation, or adding a new feature, your help is appreciated.

---

## How to contribute

1. **Fork the repository**  
   Click the "Fork" button at the top-right of the [DAI repository](https://github.com/gorankrgovic/dai).

2. **Clone your fork**  
```bash
git clone https://github.com/<your-username>/dai.git
cd dai
```

3. **Create a new branch**  
```bash
git checkout -b feature/my-new-feature
```

4. **Make changes**  
   - Follow Go best practices.
   - Format your code with `gofmt`.
   - Keep functions small and modular.

5. **Test your changes**  
```bash
go test ./...
```

6. **Commit changes**  
```bash
git add .
git commit -m "Add my new feature"
```

7. **Push to your fork**  
```bash
git push origin feature/my-new-feature
```

8. **Open a Pull Request**  
   Go to your fork on GitHub and click **"Compare & pull request"**.

---

## Code style guidelines

- Use `gofmt` before committing.
- Follow idiomatic Go naming conventions.
- Keep CLI command descriptions clear and concise.
- Avoid introducing unnecessary dependencies.

---

## Reporting issues

If you find a bug, please [open a GitHub issue](https://github.com/gorankrgovic/dai/issues) with:
- Steps to reproduce
- Expected behavior
- Actual behavior
- Any relevant logs or screenshots

---

Next: [License](license.md)
