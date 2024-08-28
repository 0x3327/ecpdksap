const EC = require("elliptic").ec;
const keccak256 = require("js-sha3").keccak256;
// Create a new elliptic curve object
const ec = new EC("secp256k1");
const { performance } = require("perf_hooks");

N = [5000, 10000, 20000, 40000, 80000];
SIMULATIONS = 3;

// Generate Ethereum Address from Public Key
function toEthAddress(PublicKey) {
  var _staX = PublicKey.getX().toString("hex");
  var _staY = PublicKey.getY().toString("hex");
  var stAA = keccak256(Buffer.from(_staX + _staY, "hex")).toString(16);
  return "0x" + stAA.slice(-40);
}

function median(values) {
  if (values.length === 0) throw new Error("No inputs");

  values.sort(function (a, b) {
    return a - b;
  });

  var half = Math.floor(values.length / 2);

  if (values.length % 2) return values[half];

  return (values[half - 1] + values[half]) / 2.0;
}

const findVariance = (arr = []) => {
  if (!arr.length) {
    return 0;
  }
  const sum = arr.reduce((acc, val) => acc + val);
  const { length: num } = arr;
  const median = sum / num;
  let variance = 0;
  arr.forEach((num) => {
    variance += (num - median) * (num - median);
  });
  variance /= num;
  return variance;
};

// Get secp256k1 generator point
const generatorPoint = ec.g;

// Get list of KeyPair
// Public Key Coordinates calculated via Elliptic Curve Multiplication
// PublicKeyCoordinates = privateKey * generatorPoint

function genDummies(N) {
  console.log("Generating ", N, " keypairs...\n");
  const PrivateKeyArray = [];
  const PublicKeyArray = [];
  for (let i = 1; i <= N; i++) {
    PrivateKeyArray.push(i.toString(16));
    PublicKeyArray.push(generatorPoint.mul(i.toString(16)));
  }
  return PublicKeyArray;
}
//console.log('PublicKeyArray:',PublicKeyArray);

// Test I
// Sender
const s = "d952fe0740d9d14011fc8ead3ab7de3c739d3aa93ce9254c10b0134d80d26a30";
const S = generatorPoint.mul(s);
const S_x = S.getX().toString("hex");
const S_y = S.getY().toString("hex");
console.log("Sender private key:\n", s + "\n------------");
console.log("Sender public  key:\n", S_x + S_y + "\n------------");

// Recipient
const p_scan =
  "0000000000000000000000000000000000000000000000000000000000000002";
const p_spend =
  "0000000000000000000000000000000000000000000000000000000000000003";
const P_scan = generatorPoint.mul(p_scan);
const P_spend = generatorPoint.mul(p_spend);
const P_scan_x = P_scan.getX().toString("hex");
const P_scan_y = P_scan.getY().toString("hex");
const P_spend_x = P_spend.getX().toString("hex");
const P_spend_y = P_spend.getY().toString("hex");
console.log("Recipient private key (scan):\n", p_scan + "\n------------");
console.log("Recipient private key (spend):\n", p_spend + "\n------------");
console.log(
  "Recipient public  key (scan):\n",
  P_scan_x + P_scan_y + "\n------------"
);
console.log(
  "Recipient public  key (spend):\n",
  P_spend_x + P_spend_y + "\n------------"
);

//console.log('S:',S);

//console.log('S_x:',S_x);
//console.log('S_y:',S_y);

// Diffie Hellman secret between sender and recipient
const Q = P_scan.mul(s);
console.log(
  "Shared Secret:\n",
  Q.getX().toString("hex") + Q.getY().toString("hex") + "\n------------"
);
console.assert(
  Q.getX().toString("hex") == S.mul(p_scan).getX().toString("hex")
);
//console.log('Q2:',Q2.getX().toString('hex')+Q.getY().toString('hex'));
const Q_x = Q.getX().toString("hex");
const Q_y = Q.getY().toString("hex");
//console.log('S_x:',Q_x);
//console.log('S_y:',Q_y);

const Qxy = Buffer.from(Q_x + Q_y, "hex");
//console.log('Qxy:',Qxy);

console.log(Qxy);
const hQ = keccak256(Qxy);
console.log("hQ", hQ);
const ViewTag = hQ.slice(0, 11);
console.log("ViewTag:\n" + ViewTag + "\n------------");
//console.log('hQ1:',hQ);

const hQG = generatorPoint.mul(hQ);
//console.log('hQG:',hQG.getX().toString('hex'));
//console.log('hQG:',hQG.getY().toString('hex'));

const stA = P_spend.add(hQG);
const stealthAddress = toEthAddress(stA);
//console.log('stA:',stA.getX().toString('hex'));
//console.log('stA:',stA.getY().toString('hex'));
console.log("stealthAddress:\n", stealthAddress + "\n------------");

simulations_wo_viewtag = [];
simulations_wo_viewtag_var = [];
simulations_w_viewtag = [];
simulations_w_viewtag_var = [];

for (let k = 0; k < N.length; k++) {
  var PublicKeyArray = genDummies(N[k]);
  // START PARSING (WITHOUT VIEWTAG)
  console.log("Start parsing ", PublicKeyArray.length, " announcements...\n");
  for (let i = 0; i < SIMULATIONS; i++) {
    var times = [];
    var startTime = performance.now();
    for (let i = 0; i < PublicKeyArray.length; i++) {
      var _Pubk = PublicKeyArray[i];
      var _q = _Pubk.mul(p_scan);
      var _hq = keccak256(
        Buffer.from(
          _q.getX().toString("hex") + _q.getY().toString("hex"),
          "hex"
        )
      );
      var _hqg = generatorPoint.mul(_hq);
      var _sta = toEthAddress(_hqg);

      if (_sta == stealthAddress) {
        console.log("Success");
      }
    }

    var _q = S.mul(p_scan);
    var _hq = keccak256(
      Buffer.from(_q.getX().toString("hex") + _q.getY().toString("hex"), "hex")
    );
    var _hqg = generatorPoint.mul(_hq);
    var _stap = P_spend.add(_hqg);
    var _sta = toEthAddress(_stap);

    var endTime = performance.now();
    var deltaTime = endTime - startTime;
    times.push(deltaTime);
    console.log(`Took ${deltaTime} milliseconds`);
  }
  if (_sta == stealthAddress) {
    console.log("Success!");
  }
  console.log("derived address:", _sta);
  console.log("stealth Address:", stealthAddress);
  console.log("----------------------------------------");
  //var timesum = 0;
  //for (let i = 0; i < times.length; i++) {
  //timesum += times[i];
  //}
  //timesum /= times.length;
  var timemed = median(times);
  var timevar = findVariance(times);
  simulations_wo_viewtag.push(timemed);
  simulations_wo_viewtag_var.push(timevar);

  console.log("Avg. execution time: ", timemed);

  // WITH VIEW TAGS
  console.log(
    "\n\nStart parsing ",
    PublicKeyArray.length,
    " announcements with ViewTags...\n"
  );
  for (let i = 0; i < SIMULATIONS; i++) {
    var times = [];
    var startTime = performance.now();
    for (let i = 0; i < PublicKeyArray.length; i++) {
      var _Pubk = PublicKeyArray[i];
      var _q = _Pubk.mul(p_scan);
      var _hq = keccak256(
        Buffer.from(
          _q.getX().toString("hex") + _q.getY().toString("hex"),
          "hex"
        )
      );
      if (_hq.slice(0, 11) != ViewTag) {
        continue;
      }
      console.log("Should not reach this line");
      var _hqg = generatorPoint.mul(_hq);
      var _stap = P_spend.add(_hqg);
      var _sta = toEthAddress(_stap);
      if (_sta == stealthAddress) {
        console.log("Stealth address derived successfully.");
      }
    }

    var _q = S.mul(p_scan);
    var _hq = keccak256(
      Buffer.from(_q.getX().toString("hex") + _q.getY().toString("hex"), "hex")
    );
    if (_hq.slice(0, 11) != ViewTag) {
      console.log("Problem with ViewTag");
    } else {
      console.log("ViewTag are matching.");
    }
    var _hqg = generatorPoint.mul(_hq);
    var _stap = P_spend.add(_hqg);
    var _sta = toEthAddress(_stap);

    var endTime = performance.now();
    var deltaTime = endTime - startTime;
    times.push(deltaTime);
    console.log(`Took ${deltaTime} milliseconds`);
  }
  if (_sta == stealthAddress) {
    console.log("Success!");
  }
  console.log("derived address:", _sta);
  console.log("stealth Address:", stealthAddress);
  //var timesum = 0;
  //for (let i = 0; i < times.length; i++) {
  //timesum += times[i];
  //}
  var timemed = median(times);
  var timevar = findVariance(times);
  simulations_w_viewtag.push(timemed);
  simulations_w_viewtag_var.push(timevar);

  console.log("Avg. execution time: ", timemed);
}
console.log("simulations_wo_viewtag: ", simulations_wo_viewtag);
//console.log("simulations_wo_viewtag var: ", simulations_wo_viewtag_var);
console.log("simulations_w_viewtag: ", simulations_w_viewtag);
//console.log("simulations_w_viewtag var: ", simulations_w_viewtag_var);

// PRIVATE KEY DERIVATION
//console.log('p_spend:',BigInt(p_spend));
//console.log('hQ:', BigInt("0x"+hQ) );
//const sta = BigInt("0x"+hQ) + BigInt(p_spend, 16);
//console.log('sta:',sta.toString(16));
//console.log('sta:',sta);

//const sta_G = generatorPoint.mul(sta.toString(16));
//console.log('sta_G:',sta_G.getX().toString('hex'));
//console.log('sta_G:',sta_G.getY().toString('hex'));
