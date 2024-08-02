import { useEffect, useState, version } from "react";

const go = new Go();

const App = () => {
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    (async () => {
      for (const el of _DB.ids) {
        await generateMetaRegistry(setLoaded);
        await delay(300);
      }

      setLoaded(true);
    })();
  }, []);

  if (loaded == false) return <div>Loading...</div>;

  return (
    <div className="App">
      <div className="BackgroundImage"></div>
      <Navbar></Navbar>
      <SenderSide></SenderSide>
      <RecipientSide></RecipientSide>
      <TxList></TxList>
    </div>
  );
};

let txListUpdate;

const Navbar = () => {
  return (
    <div className="Navbar">
      <h1 className="Contrast">ECPDKSAP</h1>
    </div>
  );
};

const SenderSide = () => {
  const [selectedMetaAddress, setSelectedMetaAddress] = useState(
    _DB.metaAddrRegistry[0]
  );

  const [senderInfo, setSenderInfo] = useState({
    privKey: ``,
    pubKey: ``,
  });

  const [stealthInfo, setStealhInfo] = useState({
    publicKey: ``,
    address: ``,
    viewTag: ``,
  });

  useEffect(() => {}, []);

  return (
    <div className="SenderSide">
      <div className="MetaAddressRegistry">
        <h3 className="Header">Meta Registry</h3>

        <div className="MetaIdsWrapper">
          {_DB.metaAddrRegistry.map((e, k) => {
            console.log(selectedMetaAddress.id);
            const _className =
              selectedMetaAddress.id === e.id
                ? "MetaId SelectedMetaId"
                : "MetaId";
            return (
              <div
                className={_className}
                key={k}
                onClick={async () => {
                  setSelectedMetaAddress(e);
                  await delay(100);
                  await calculateStealth(
                    setSenderInfo,
                    {
                      K: e.K,
                      V: e.V,
                      Version: "v2",
                      ViewTagVersion: "v0-1byte",
                    },
                    setStealhInfo,
                    true
                  );
                }}
              >
                {e.id}
              </div>
            );
          })}
        </div>

        <div className="MetaPublicKeys">
          {/* <h4>Public keys: </h4> */}
          <div className="Entry">
            <div className="Label">SpendingPubKey:</div>
            <div className="Value">{selectedMetaAddress.K}</div>
          </div>
          <div className="Entry">
            <div className="Label">ViewingPubKey:</div>
            <div className="Value">{selectedMetaAddress.V}</div>
          </div>
        </div>

        <div className="SenderInfo">
          <div className="Header">
            <h3>Sender info:</h3>
            <button
              onClick={() => {
                calculateStealth(
                  setSenderInfo,
                  {
                    K: selectedMetaAddress.K,
                    V: selectedMetaAddress.V,
                    Version: "v2",
                    ViewTagVersion: "v0-1byte",
                  },
                  setStealhInfo
                );
              }}
            >
              Generate
            </button>
          </div>

          <div className="Body">
            <div className="Entry">
              <div className="Label">SenderPrivKey:</div>
              <div className="Value">{senderInfo.privKey}</div>
            </div>

            <div className="Entry">
              <div className="Label">SenderPubKey:</div>
              <div className="Value">{senderInfo.pubKey}</div>
            </div>
          </div>
        </div>

        <div className="StealthInfo">
          <div className="Header">
            <h3>Stealth info:</h3>
            <button
              onClick={() => {
                const info = [];
                txListUpdate([
                  {
                    senderPubKey: senderInfo.pubKey,
                    viewTag: stealthInfo.viewTag,
                  },
                ]);
              }}
            >
              Send ETH
            </button>
          </div>

          <div className="Body">
            <div className="Entry">
              <div className="Label">Public key:</div>
              <div className="Value">{stealthInfo.publicKey}</div>
            </div>

            <div className="Entry">
              <div className="Label">Address:</div>
              <div className="Value">{stealthInfo.address}</div>
            </div>

            <div className="Entry">
              <div className="Label">ViewTag:</div>
              <div className="Value">{stealthInfo.viewTag}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const RecipientSide = () => {
  const [selectedMetaAddress, setSelectedMetaAddress] = useState(
    _DB.metaAddrRegistry[0]
  );

  const [discoveredData, setDiscoveredStealthInfo] = useState([]);

  const updateDiscoveredStealthInfo = (update) => {
    //   const onlyUnique = (value, index, array) => {
    //     return array.indexOf(value) === index;
    //   };

    //   setDiscoveredStealthInfo([...discoveredData, ...update].filter(onlyUnique));

    setDiscoveredStealthInfo([...update]);
  };

  useEffect(() => {}, []);

  return (
    <div className="RecipientSide">
      <div className="MetaAddressRegistry">
        <h3 className="Header">Meta Registry</h3>

        <div className="MetaIdsWrapper">
          {_DB.metaAddrRegistry.map((e, k) => {
            console.log(selectedMetaAddress.id);
            const _className =
              selectedMetaAddress.id === e.id
                ? "MetaId SelectedMetaId"
                : "MetaId";
            return (
              <div
                className={_className}
                key={k}
                onClick={() => {
                  setSelectedMetaAddress(e);
                }}
              >
                {e.id}
              </div>
            );
          })}
        </div>

        <div className="MetaPublicKeys">
          {/* <h4>Public keys: </h4> */}
          <div className="Entry">
            <div className="Label">SpendingPubKey:</div>
            <div className="Value">{selectedMetaAddress.K}</div>
          </div>
          <div className="Entry">
            <div className="Label">ViewingPubKey:</div>
            <div className="Value">{selectedMetaAddress.V}</div>
          </div>
        </div>

        <div className="SenderInfo">
          <div className="Header">
            <h3>Recipient info:</h3>
            <button
              onClick={() => {
                recieveScan(selectedMetaAddress, updateDiscoveredStealthInfo);
              }}
            >
              Scan
            </button>
          </div>

          <div className="Body">
            <div className="Entry">
              <div className="Label">SpendingPrivKey:</div>
              <div className="Value">{selectedMetaAddress.k}</div>
            </div>

            <div className="Entry">
              <div className="Label">ViewingPrivKey:</div>
              <div className="Value">{selectedMetaAddress.v}</div>
            </div>
          </div>
        </div>

        <div className="DiscoveredStealthInfo">
          <div className="Header">
            <h3>Discovered Stealth info:</h3>
          </div>
          {discoveredData.map((el, k) => {
            return (
              <div className="Row" key={k}>
                <div className="Body">
                  <div className="Entry">
                    <div className="Label">Address:</div>
                    <div className="Value">{el.address}</div>
                  </div>
                  <div className="Entry">
                    <div className="Label">Priv. Key:</div>
                    <div className="Value">{el.privKey}</div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

const TxList = () => {
  const [data, _txListUpdate] = useState([]);
  txListUpdate = (update) => {
    console.log({ update });
    _DB.txList = [...data, ...update];
    _txListUpdate(_DB.txList);
  };
  return (
    <div className="TxList">
      <h3 className="Header">(Public) Transaction list:</h3>
      <div className="Table">
        <div className="Head">
          <div className="Row">
            <div className="Entry">#</div>
            <div className="Entry">Sender Pub. Key</div>
            <div className="Entry">View Tag</div>
          </div>
        </div>
        <div className="Body">
          {data.map((el, k) => {
            return (
              <div className="Row" key={k}>
                <div className="Entry">{k + 3327}</div>
                <div className="Entry">{el.senderPubKey}</div>
                <div className="Entry">{el.viewTag}</div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default App;

const _DB = {
  ids: ["Marija", "Mihajlo", "Mihailo", "Alex", "Malisa", "Milos"],
  idProccessed: -1,
  metaAddrRegistry: [],

  txList: [{ senderPubKey: "", viewTag: "" }],
};

const calculateStealth = async (
  setSenderInfo,
  info,
  setStealhInfo,
  shortcircuit
) => {
  if (shortcircuit == null) {
    WebAssembly.instantiateStreaming(
      fetch("ecpdksap.wasm"),
      go.importObject
    ).then((result) => {
      go.argv = ["js", "gen-send-info"];

      go.run(result.instance);

      const info2 = JSON.parse(window.senderMeta);

      console.log("generateSenderInfo", { info2 });

      setSenderInfo({ privKey: info2.r, pubKey: info2.R });
    });
  }

  await delay(100);
  generateStealth(info, setStealhInfo);

  // WebAssembly.instantiateStreaming(
  //   fetch("ecpdksap.wasm"),
  //   go.importObject
  // ).then((result) => {
  //   const info2 = JSON.parse(window.senderMeta);

  //   console.log("generateSenderInfo", { info2 });

  //   info.r = info2.r;

  //   go.argv = ["js", "send", JSON.stringify(info)];
  //   go.run(result.instance);

  //   setStealhInfo({
  //     publicKey: window.StealthPubKey,
  //     address: window.StealthAddress,
  //     viewTag: window.StealthViewTag,
  //   });
  // });
};

const generateStealth = (info, setStealhInfo) => {
  WebAssembly.instantiateStreaming(
    fetch("ecpdksap.wasm"),
    go.importObject
  ).then((result) => {
    const info2 = JSON.parse(window.senderMeta);

    console.log("generateStealth", { info2 });

    info.r = info2.r;

    go.argv = ["js", "send", JSON.stringify(info)];
    go.run(result.instance);

    setStealhInfo({
      publicKey: window.StealthPubKey,
      address: window.StealthAddress,
      viewTag: window.StealthViewTag,
    });
  });
};

const recieveScan = (info, updateDiscoveredStealthInfo) => {
  WebAssembly.instantiateStreaming(
    fetch("ecpdksap.wasm"),
    go.importObject
  ).then((result) => {
    const Rs = _DB.txList.map((el) => el.senderPubKey);
    const ViewTags = _DB.txList.map((el) => el.viewTag);

    console.log({ info });

    go.argv = [
      "js",
      "receive-scan",
      JSON.stringify({
        k: info.k,
        v: info.v,
        Version: "v2",
        ViewTagVersion: "v0-1byte",
        Rs,
        ViewTags,
      }),
    ];
    go.run(result.instance);

    const dStealthAddrs = window.DiscoveredStealthAddrs.split(".");
    const dPrivKeys = window.DiscoveredStealthPrivKeys.split(".");

    const discovered = dStealthAddrs.map((address, k) => {
      return {
        address,
        privKey: dPrivKeys[k],
      };
    });

    updateDiscoveredStealthInfo(discovered);
  });
};

const delay = (ms) => new Promise((res) => setTimeout(res, ms));

const generateMetaRegistry = async (setLoaded) => {
  WebAssembly.instantiateStreaming(
    fetch("ecpdksap.wasm"),
    go.importObject
  ).then((result) => {
    go.argv = ["js", "gen-recipient-info"];

    go.run(result.instance);

    const info = JSON.parse(window.recipientMeta);

    console.log("generateMetaRegistry", { info });

    _DB.idProccessed += 1;

    _DB.metaAddrRegistry.push({ id: _DB.ids[_DB.idProccessed], ...info });
  });
};
