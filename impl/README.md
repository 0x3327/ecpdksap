## Benchmark

From `./benchmark` run:

```
go test . -bench=<Specific-Benchmark> -benchmem -benchtime=1x -timeout 2000m > log_results.txt
```

i.e:

```
go test . -bench=Benchmark_Curves_100 -benchmem -benchtime=1x -timeout 2000m > log_results.txt
```

## Example CLI inputs / cmds

### Generation of example

```
go run . gen-example <version: v0 | v1 | v2> <sample-size: 1...1000>
```

### Send

```
export SEND_EXAMPLE_INPUT=$(cat ./gen_example/example/inputs/send.json) && go run . send $SEND_EXAMPLE_INPUT
```

### Receive

```
export RECEIVE_EXAMPLE_INPUT=$(cat ./gen_example/example/inputs/receive.json) && go run . receive-scan $RECEIVE_EXAMPLE_INPUT
```
