# Reporting Website Changes

Miru is designed to allow for as much flexibility as possible for administrator's to check websites
for changes via scripts that are periodically run by miru itself.  These scripts are expected to
follow a couple of conventions to work as seamlessly as possible.  These conventions are as follows:

1. Scripts must read their input from `stdin`.
2. Scripts should only write to `stdout` to output report [JSON](https://en.wikipedia.org/wiki/JSON#Example).

More on these points after a brief explanation of how miru runs scripts.

## Report JSON

To make communication between monitor scripts and miru as simple as possible, the only mechanism
that is used is the operating system's standard input and standard output, referred to as
`stdin` and `stdout` respectively.

When a script is run by miru, miru will write the last report produced by the script to the
script's `stdin`.  This allows your script to make more informed decisions about how to evaluate
a website (i.e. check for changes) with reference to any information that was gathered about the
site on the last run.

In order to create a new report, which will become the input written to the script the next time
it is run, your script only has to write a new report to its `stdout`.  This does mean that your
script is a little restricted.  If you would like to write more information to the terminal for
debugging purposes, you will have to write it to `stderr` (standard error).

### Reading from stdin

The following examples show how you can quickly and easily read from `stdin` in each supported
language and parse the [JSON](https://en.wikipedia.org/wiki/JSON#Example) into a data structure that your program can work with.

**Python:**

```python
import json, sys

report = json.load(sys.stdin)
# You can now use report as a regular Python dictionary!
```

**Ruby:**

```ruby
require 'JSON'

report = JSON.load STDIN.read
# You can now use report as a regular Ruby hash!
```

### Writing to stdout

Writing to `stdout` should be a lot more familiar.  Your preferred scripting language's
equivalent of `print` or `puts` will work!

**Python:**

```python
print({"message": "hello from Python! :)"})
```

**Ruby:**

```ruby
puts {"message" => "hello from Ruby! :)"}
```

### Debugging Messages

In case your script needs to output some data for debugging purposes, `stderr` must be used
to avoid getting conflicted with the [JSON](https://en.wikipedia.org/wiki/JSON#Example) output expected by miru.

**Python:**

```python
import sys

sys.stderr.write("Hello from Python's stderr!")
```

**Ruby:**

```ruby
STDERR.write "Hello from Ruby's stderr!"
```

## Report Format

The reports that monitor scripts are expected to write are essentially just a
[JSON](https://en.wikipedia.org/wiki/JSON#Example) object with a handful of expected values.
The format is as follows:

```json
{
  "lastChangeSignificance": 0,
  "message": "Please investigate this site",
  "checksum": "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03",
  "state": {}

}
```

### Change Significance

The `lastChangeSignificance` field is a measure how important a change to a site is, and is measured
with a number between 0 and 4 inclusive.  The five levels are as follows:

| Level | Name             | Description                                          |
|-------|------------------|------------------------------------------------------|
| 0     | `no_change`      | Nothing on the page has changed.                     |
| 1     | `minor_update`   | A minor textual change occurred, such as a typo fix. |
| 2     | `content_change` | The content has been modified in a meaningful way.   |
| 3     | `rewritten`      | The entire content has been replaced.                |
| 4     | `deleted`        | The site has been completely deleted.                |

### Message

The `message` field can actually contain any message you'd like administrators to see in miru.

### Checksum

The `checksum` field should be a [SHA256 hash](https://en.wikipedia.org/wiki/SHA-2) of whatever site content
that your script is reading to determine whether a change took place.  In a nutshell, this hash value will
change whenever anything in the text used to compute it changes. Two inputs with even a minor difference will
hash to completely different values, giving you a quick way to do a really simple test for changes.

There are builtin libraries in Python and Ruby that help you compute these hashes.

**Python:**

```python
import hashlib

hasher = hashlib.sha256().update('any string can go here')
checksum = hasher.hexdigest()
```

**Ruby:**

```ruby
require 'digest'

checksum = Digest::SHA256.hexdigest 'any string can go here'
```

### State

The `state` field can be any JSON object you'd like. This field is present to allow your script flexibility to
include any extra data that you would to store so that it becomes input to your script on successive runs.
