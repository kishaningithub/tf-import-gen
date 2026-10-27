# Agent instructions

## Tooling

- If [pkgx](https://pkgx.sh) is available on the machine, use it to run any tools you need (e.g. `pkgx zstd …`, `pkgx python@3.12 …`) instead of installing them globally.

## Code style

- Ensure the code is very readable. Readability is greater than performance.
- Use [samber/lo](https://github.com/samber/lo) wherever applicable to make the code more readable.

## Tests

- Use the `require` package of [stretchr/testify](https://github.com/stretchr/testify) to keep the tests readable and free of boilerplate.
