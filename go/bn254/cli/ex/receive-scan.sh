set -e

export RECEIVE_EXAMPLE_INPUT=$(cat ./cli/ex/inputs/receive.json) && go run . receive-scan $RECEIVE_EXAMPLE_INPUT
    