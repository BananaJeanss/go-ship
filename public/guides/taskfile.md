# Taskfile

<img src="/public/guides/assets/Task.svg" alt="Taskfile Logo" width="250">

Taskfile is a modern build tool/Make alternative, written in Go, which lets you write build scripts in YAML and execute them via `task`.

This can be incredibly useful for repetitive commands, specific build configs, scripts, basically anything that can be run via the command line.

Taskfile runs on all platforms, has an easy setup, and a ton of other features.

Though optional, it's a really nice quality-of-life tool to have in your Go projects, such as combining multiple commands, repetitive docker deployments, specific build commands, and more.

## Installation

Follow the OS specific instructions on the [Taskfile installation guide](https://taskfile.dev/docs/installation).

## Usage

With taskfile installed, cd into your project directory and run `task --init` to create a new `Taskfile.yml` in your project.

This will include the default Taskfile.yml contents, which should look like this:

```yaml
# yaml-language-server: $schema=https://taskfile.dev/schema.json

version: "3"

vars:
  GREETING: Hello, world!

tasks:
  default:
    desc: Print a greeting message
    cmds:
      - echo "{{.GREETING}}"
    silent: false
```

---

Let's break it down.

- Version: The taskfile schema version, which should be left as is.
- Vars: Variables that can be used in your tasks. `GREETING` is set to "Hello, world!", so when you use `{{.GREETING}}` in your tasks, it will be replaced with "Hello, world!".<br>This uses Go's `text/template` syntax, so you can use any of the features of Go templates in your taskfile.
- Tasks: The actual tasks you want to run.
  - The top level key is the task name, in this case `default`. `default` is a special task name that will be run when you run `task` without any arguments.
  - `desc`: A description of what the task does, which will be shown when you run `task --list`.
  - `cmds`: The commands to run when the task is executed. In this case, it will echo the value of `GREETING`, which is "Hello, world!".
  - `silent`: Hides taskfile's echo (excluding errors), but does not hide the output of the commands themselves.

Now, if you run `task` in your terminal, it will execute the `default` task and print:

```bash
$ task
task: [default] echo "Hello, world!"
Hello, world!
```

But how about we make it do something useful for our Go project?

---

Let's write a task that tidies up our Go imports and code.

```yaml
tidy:
  desc: Tidy up Go imports and code
  cmds:
    - go mod tidy
    - gofmt -w ./...
  silent: true
```

This will execute `go mod tidy` to clean up unused dependencies and download missing ones, and `gofmt -w ./...` to format all Go files in the current directory and subdirectories.

Let's move on to more advanced examples.

---

Let's say we want to include the git commit hash in our build. We can do this by adding a build task to the Taskfile.yml:

```yaml
build:
  desc: Builds the Go project with the git commit hash
  cmds:
    - go build -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)" -o dist/ .
  silent: true
```

For this to work for your project, you will need to make sure the package named `main` has a `commitHash` variable defined (`var commitHash string`), or you can change the `-X main.commitHash` part to match the package and variable name you want to use.

Now, when run, the -X flag changes the value of the `commitHash` variable in your Go code to the output of `git rev-parse --short HEAD`, and `-o dist/ .` builds the project and outputs the binary to the `dist` directory. Pretty cool, right?

---

But what if we have a .env variable that we want to use in our command?
No problem, we can use `dotenv` to load the .env variables into our Taskfile.yml.

Add this to the top of your Taskfile.yml:

```yaml
dotenv: [".env"]
```

This will load variables from .env and any other file you specify in the array.

Let's say we have a script that needs to execute something on a database.
We can use the environment variables loaded from .env in the command like this:

```yaml
migrate-up:
  desc: Apply all pending DB migrations
  cmds:
    - goose -dir db/migrations postgres "$DB_URL" up
  silent: true
```

`$DB_URL` will be replaced with the value of `DB_URL` from the .env file, and the command will be executed with that value.

Of course, `$` will also take system environment variables.

---

These are just a few examples of what you can do with Taskfile, and should be enough for getting started. Obviously, Taskfile can be used outside of Go as well, and has a ton of other features, such as task dependencies, arguments, conditionals, and more.

You can read more about these features in the official [Taskfile documentation](https://taskfile.dev/docs/guide).
