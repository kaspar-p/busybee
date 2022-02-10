# Testing

The driving philosophy behind testing busybee is that everything should be tested, preferably twice. The added confidence of being able to run tests and just _know_ that everything is working correctly is the goal.

Busybee is currently tested in two different environments, local machine and Github Actions. This manifests only in where the tests get their environment variables.

The tests use the default Golang `testing` library, as well as a discord integration testing library [Gourd](https://github.com/kaspar-p/gourd).

There is a separate bot called `busybee-test` that the integration tests are run against.

## Local Machine

To run the tests locally, install the project and run `go test ./...` from the root of the project. This finds every test within the project. The environment variables for local machine tests are kept in a `.gitignore`'d file `.env.test.yaml` at the root level. The structure of this file follows the `.env.template.yaml` structure, ALONG WITH an added structure

```
GOURD_BOT:
    APP_ID:
    TOKEN:
```

This information is required for the integration tests to work.

## Github Actions

These tests are run on the code in pull requests made to the `master` branch. These are run automatically on push to the external branch. The tests rely on Github Secrets for their environment variables. The naming of the variables are such that the code does not have to change. That is, the nested `.yaml` structure:

```
LEVEL_ONE:
    LEVEL_TWO:
        LEVEL_THREE: some_value
```

Yields a key of `LEVEL_ONE.LEVEL_TWO.LEVEL_THREE`. The Github secrets a flat key/value map, but the names of the keys stay the same:

```
LEVEL_ONE.LEVEL_TWO.LEVEL_THREE: some_value
```

This is identical to how environment variables are stored for production.
