set -e

export SEND_EXAMPLE_INPUT=$(cat ./cli/ex/inputs/send.json) && go run . send $SEND_EXAMPLE_INPUT
