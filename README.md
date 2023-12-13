# aws-swf-er-poc

## Development

### Requirements

- Go 1.19+
- GNU Make
- git

### Build & Execute

When running on a macos/arm host:

```bash
make clean compile-only
./bin/darwin/worker
```

To compile, run tests and linter
```bash
make
```
