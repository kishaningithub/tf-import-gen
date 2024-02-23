# Contributing

## Setting up the code in local

- Fork and clone this repo.
- Ensure you have the latest version of [golang](https://go.dev/) and [make](https://www.gnu.org/software/make/)
  installed.
- To format the code, run command `make fmt`.
- To build the project, run command `make build`.
- To run only the tests, run command `make test`.

## Upgrading dependencies
- For code, run command `make update-deps`
- For actions, go inside the yamls present in the [workflows](.github/workflows/) and bump the versions.