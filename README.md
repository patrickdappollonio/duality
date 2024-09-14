# `duality`: run multiple commands in parallel

`duality` is a Go program that allows you to run multiple commands in parallel and see their output interleaved in real-time.

When `duality` starts, each command gets a unique ID and its output is prefixed with that ID plus whether the output was rendered to `stdout` or `stderr` in the format `[ID][stdout|stderr] log message`.

Each command is run in parallel and their output is interleaved in real-time.  If a command exits, `duality` will print the exit status and the command's ID. If the exit code is non-zero, `duality` will exit with a non-zero exit code.

If any of the command fails, all other commands are sent a `SIGTERM` signal, but `duality` won't wait for them to be closed, then `duality` will exit with a non-zero exit code.

If `duality` receives a `SIGINT` or `SIGTERM` signal, it will forward a `SIGTERM` signal to all commands and exit with a non-zero exit code. Killing `duailty` will also kill all commands.

By default, all commands are run through a shell, using `/bin/bash -c` as a default. This will allow you to do some more advanced commands, like `sleep 10 && echo "Hello"`. You can change the shell used by setting the `DUALITY_SHELL_COMMAND` environment variable (e.g. `export DUALITY_SHELL_COMMAND="/bin/ash -c"`).

### Download

You can download the latest release from the [releases page](https://github.com/patrickdappollonio/duality/releases).

### Usage

Simply run the command with arguments, where each argument is a single command. If you want to ensure the commands are processed appropriately, ensure they're quoted properly.

```bash
$ duality "echo Hello" "echo World"
[75ra32uh] running: /bin/bash -c "echo World"
[52updl62] running: /bin/bash -c "echo Hello"
[52updl62][stdout] Hello
[52updl62] "echo Hello" executed successfully (exit code: 0)
[75ra32uh][stdout] World
[75ra32uh] "echo World" executed successfully (exit code: 0)
```

The log output will tell you what command and what shell interpreter will be used to run the command. Any log statements sent to `stdout` or `stderr` will be prefixed with the command ID and whether it was sent to `stdout` or `stderr`. When a command finishes, `duality` will print the exit status and the command ID. On any failed command, `duality` will exit with a non-zero exit code and report the error back to you, without prefixes.

At least one command must be provided, otherwise `duality` will exit with a non-zero exit code. Providing empty commands in between non-empty commands is allowed, but the empty commands will be skipped.

### Usage with Docker

In Docker, we provide a convenience image that you can use to run `duality` in a container:

```sh
ghcr.io/patrickdappollonio/duality:v1
```

Note this image is based on busybox, and as such, the shell used is `sh`. The binary is located at `/usr/local/bin/duality`.

If you need a different shell, you can also use the image `ghcr.io/patrickdappollonio/duality` as a build step to easily download the binary, then copy it to your own image, effectively allowing you to change to any shell you want. Here's an example:

```Dockerfile
FROM ghcr.io/patrickdappollonio/duality:v1 as duality

FROM bash:latest
COPY --from=duality /usr/local/bin/duality /usr/local/bin/duality
```

By having it in your own image, you set the rules: you can use any shell (just make sure you tell `duality` which shell to use by setting the `DUALITY_SHELL_COMMAND` environment variable) and you can change entrypoints and other settings as you see fit.

### About IDs

The IDs are generated based on the command being executed and its position as provided to the `duality` program. In practice, a command will always yield the same ID as long as:

* The string representation of the command hasn't changed
* The position of the command, as provided to the CLI, hasn't changed

This process makes IDs predictable and easy to `grep`.

If you want to maintain a position after deleting a command, you can pass an empty string as a parameter. This will ensure the ID is kept, but the empty command will be skipped:

```bash
$ duality "foo" "bar" "baz"
[9zmorkes] running: /bin/bash -c "baz"
[qex9920b] running: /bin/bash -c "foo"
[2mlvq8wb] running: /bin/bash -c "bar"

$ duality "foo" "" "baz"
[9zmorkes] running: /bin/bash -c "baz"
[qex9920b] running: /bin/bash -c "foo"
```

Note in the example above, the IDs `qex9920b` and `9zmorkes` were maintained, even when now the second command is empty.
