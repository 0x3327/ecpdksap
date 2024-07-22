const go = new Go();

WebAssembly.instantiateStreaming(fetch("test-sap.wasm"), go.importObject).then(
  (result) => {
    console.log({ x: go.argv });
    go.argv.push("send");
    go.run(result.instance);
    // document.getElementById("log").innerText = `${MainSend(123)}`;
  }
);
