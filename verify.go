package main

import "github.com/piotrnar/gocoin/lib/secp256k1"
import "encoding/hex"
import "encoding/base64"
import "crypto/sha256"
import "os"

func verify(pubkey, signature, message string) {
  mlen := len(message)
  if mlen >= 0xfd { os.Exit(1) } // TODO not supported
  btcmsg := "\x18Bitcoin Signed Message:\n"
  btcmsg += string([]byte{byte(mlen)})
  btcmsg += message
  msghash1 := sha256.Sum256([]byte(btcmsg))
  msghash2 := sha256.Sum256(msghash1[:])
  sigraw, _ := base64.StdEncoding.DecodeString(signature)
  var sig secp256k1.Signature
  // TODO check recId, 27 and compressed from sigraw[0]
        sig.R.SetBytes(sigraw[1:1+32])
        sig.S.SetBytes(sigraw[1+32:1+32+32])
  var pubkey_xy secp256k1.XY
  xy, _ := hex.DecodeString(pubkey)
  pubkey_xy.ParsePubkey(xy)
  var msg secp256k1.Number
  msg.SetBytes(msghash2[:])
  if !sig.Verify(&pubkey_xy, &msg) { println("failed") }
}

func main() {
  // addr := "1QHBj5GjAEp7oFKhp5QdeLXW8jnm2PupBs"
  pubkey := "04a4c6b2c8c62ffe17c6740434776c912ff3fc11891922fac4ee93e966817f4446caa82ba4fe6b494f67f870c332fdc5ba22fcdfd6b3a1e88e8bb0693b9bf2c1c2"
  signature := "HKnOcPe/RxF48z5U6JbyetZC7+wmPrlUOumbbecpMVlwbcfGLlTwGBtMDzjD4wxOg/VjQDg7TxHqP/Mfoohp7Cs="
  message := "test"
  verify(pubkey, signature, message)
}
