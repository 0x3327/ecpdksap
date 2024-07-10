## Example inputs / cmds

### Send

```
export SEND_EXAMPLE_INPUT=$(cat ./cli/ex/inputs/send.json) && go run . send $SEND_EXAMPLE_INPUT
```

### Receive (without using View Tag)

```
export RECEIVE_EXAMPLE_INPUT=$(cat ./cli/ex/inputs/receive.json) && go run . receive-scan $RECEIVE_EXAMPLE_INPUT
```

### Receive (usign View Tag)

```
export RECEIVE_EXAMPLE_INPUT=$(cat ./cli/ex/inputs/receive-using-vtag.json) && go run . receive-scan-using-vtag $RECEIVE_EXAMPLE_INPUT
```
