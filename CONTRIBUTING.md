# Contributing

## Welcome

Contributions are very welcome on go-api-schemas. When contributing please keep this in mind:

- Open an issue to discuss new bigger features.
- Write code consistent with the project style and make sure the tests are passing.
- Stay in touch with us if we have follow up questions or requests for further changes.

## Development

### Local Environment

- Set `AIVEN_TOKEN` and `AIVEN_PROJECT_NAME` environment variables to your Aiven API token and project name.

- Run `task run` to generate and persist schemas.

You can also use `task build` to build the binary which you can then run the binary with `./go-api-schemas`.

Try it and see `./go-api-schemas --help` for more configuration options besides the environment variables.

### Tests

```bash
task test
```

### Static checking and Linting

We use [Trunk.io](https://trunk.io/) for static checking and linting. Install it locally and you'll be good to go.

## Opening a PR

- Commit messages should describe the changes, not the filenames. Win our admiration by following the
  [excellent advice from Chris Beams](https://chris.beams.io/posts/git-commit/) when composing commit messages.
- Choose a meaningful title for your pull request.
- The pull request description should focus on what changed and why.
- Check that the tests pass (and add test coverage for your changes if appropriate).

### Commit Messages

This project adheres to the [Conventional Commits](https://conventionalcommits.org/en/v1.0.0/) specification.
Please, make sure that your commit messages follow that specification.
