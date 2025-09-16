# Contributing

## Welcome

Contributions are very welcome to go-api-schemas. When contributing, please keep the following in mind:

- Open an issue to discuss new bigger features.
- Write code consistent with the project style and make sure the tests are passing.
- Stay in touch with us if we have follow up questions or requests for further changes.

## Development

During development, you might want to see how schemas are generated,
or use them with the [Aiven Terraform Provider][Aiven Terraform Provider] or [Aiven Kubernetes Operator][Aiven Kubernetes Operator].

**Note:** Generated files should not be pushed as part of a PR unless specifically requested.

### Generating schemas

To generate schemas with the default URL, run:

```shell
AIVEN_PROJECT_NAME=my-project AIVEN_TOKEN=abc123 task run
```

Where `AIVEN_PROJECT_NAME` is the name of your project, and `AIVEN_TOKEN` is your personal token.

To generate schemas with a custom URL, run:

```shell
AIVEN_WEB_URL=https://custom-api AIVEN_PROJECT_NAME=my-project AIVEN_TOKEN=abc123 task run
```

This will:

- Download JSON schema files from the Aiven API
- Generate new schema files in the `pkg/dist` directory

### Using the generated schemas

You can make the Terraform Provider or Kubernetes Operator use the generated `go-api-schemas` repository for development.
You can test your changes even before committing them.

1. Clone the [Aiven Terraform Provider][Aiven Terraform Provider] or [Aiven Kubernetes Operator][Aiven Kubernetes Operator].
2. Within the cloned repository, create a `go.work` file with the following content:

   ```go.work
   go 1.24.0

   use (
       .
       ../your-go-api-schemas-repository-path
   )
   ```

   The `go 1.24.0` version must match the version in the `go.mod` file,
   which may have been updated since these instructions were written.

3. Run the following command in the Terraform or Operator directory:

   ```shell
   task generate
   ```

   Your changes will now be visible in the [Aiven Terraform Provider][Aiven Terraform Provider] or [Aiven Kubernetes Operator][Aiven Kubernetes Operator] code.
   Follow each project's documentation for next steps.

### Advanced generation options

For more options, run:

```shell
go run ./... --help
```

### Tests

```bash
task test
```

### Static checking and linting

We use [Trunk.io](https://trunk.io/) for static checking and linting. Install it locally and you'll be ready to go.

## Opening a PR

- Commit messages should describe the changes, not the filenames. Win our admiration by following the
  [excellent advice from Chris Beams](https://chris.beams.io/posts/git-commit/) when composing commit messages.
- Choose a meaningful title for your pull request.
- The pull request description should focus on what changed and why.
- Check that the tests pass (and add test coverage for your changes if appropriate).

### Commit Messages

This project adheres to the [Conventional Commits](https://conventionalcommits.org/en/v1.0.0/) specification.
Please ensure that your commit messages follow this specification.

[Aiven Terraform Provider]: https://github.com/aiven/terraform-provider-aiven
[Aiven Kubernetes Operator]: https://github.com/aiven/aiven-operator
