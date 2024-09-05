cp main.go ./builds/tmp-main.mgo

export MAIN_GO_TMP=$(cat main.go) && \
grep -l $MAIN_GO_TMP main.go | sort | uniq | xargs perl -e "s/\/\/wasm.build::://" -pi

GOOS=js GOARCH=wasm go build -o  ./builds/ecpdksap.wasm

mv ./builds/tmp-main.mgo main.go 