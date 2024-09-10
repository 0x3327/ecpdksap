## ECPDKSAP Build process

### Standard Binary Files

**AMD**

```bash
GOARCH=amd64 go build -o ./builds/ecpdksap-amd
```

**ARM**

```bash
GOARCH=arm64 go build -o ./builds/ecpdksap-arm
```

### Building a WASM file

```bash
cp main.go ./builds/tmp-main.mgo && \
export MAIN_GO_TMP=$(cat main.go) && \
grep -l $MAIN_GO_TMP main.go | sort | uniq | xargs perl -e "s/\/\ wasm.build::://" -pi && \
GOOS=js GOARCH=wasm go build -o  ./builds/ecpdksap.wasm && \
mv ./builds/tmp-main.mgo main.go
```
