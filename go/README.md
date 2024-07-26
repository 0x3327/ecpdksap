## Benchmark

From `./benchmark` run:

```
go test -v . -bench=. -benchtime=1x
```

## Example CLI inputs / cmds

### Generation of example

```
go run . gen-example <version: v0 | v2> <sample-size: 1...1000>
```

### Send

```
export SEND_EXAMPLE_INPUT=$(cat ./gen_example/example/inputs/send.json) && go run . send $SEND_EXAMPLE_INPUT
```

### Receive (without using View Tag)

```
export RECEIVE_EXAMPLE_INPUT=$(cat ./gen_example/example/inputs/receive.json) && go run . receive-scan $RECEIVE_EXAMPLE_INPUT
```
