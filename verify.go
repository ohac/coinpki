package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/itchyny/base58-go"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	"os"
)

func verify(pubkey, signature, message string) {
	mlen := len(message)
	if mlen >= 0xfd { // TODO not supported
		os.Exit(1)
	}
	btcmsg := "\x18Bitcoin Signed Message:\n"
	btcmsg += string([]byte{byte(mlen)})
	btcmsg += message
	msghash1 := sha256.Sum256([]byte(btcmsg))
	msghash2 := sha256.Sum256(msghash1[:])
	sigraw, _ := base64.StdEncoding.DecodeString(signature)
	var sig Signature
	r0 := sigraw[0] - 27
	recid := int(r0 & 3)
	compressed := (r0 & 4) == 1
	fmt.Println("sigraw[0]:", sigraw[0], recid, compressed)
	if compressed { // TODO not supported
		os.Exit(1)
	}
	sig.R.SetBytes(sigraw[1 : 1+32])
	sig.S.SetBytes(sigraw[1+32 : 1+32+32])
	var pubkey_xy XY
	xy, _ := hex.DecodeString(pubkey)
	pubkey_xy.ParsePubkey(xy)
	var msg Number
	msg.SetBytes(msghash2[:])
	if !sig.Verify(&pubkey_xy, &msg) {
		println("failed")
	}

	pubkey_xy.Print("pub1")

	var pubkey_xy2 XY

	ret2 := sig.recover(&pubkey_xy2, &msg, recid)
	if !ret2 {
		println("failed")
	}
	pubkey_xy2.Print("pub2")

	out := make([]byte, 65)
	pubkey_xy2.GetPublicKey(out)
	for _, x := range out {
		fmt.Printf("%02x", x)
	}
	fmt.Println()
	fmt.Println(pubkey)

	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(out)
	pub_hash_1 := sha256_h.Sum(nil)

	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)

	bcpy := append([]byte{byte(0)}, pub_hash_2...)
	sha256_h.Reset()
	sha256_h.Write(bcpy)
	hash1 := sha256_h.Sum(nil)
	sha256_h.Reset()
	sha256_h.Write(hash1)
	hash2 := sha256_h.Sum(nil)
	bcpy = append(bcpy, hash2[0:4]...)

	z := new(big.Int)
	z.SetBytes(bcpy)
	zstr := z.String()

	enc := base58.BitcoinEncoding
	encdd, _ := enc.Encode([]byte(zstr))

	s := string(encdd)
	for _, v := range bcpy {
		if v != 0 {
			break
		}
		s = "1" + s
	}

	fmt.Println(s)
	fmt.Println("1QHBj5GjAEp7oFKhp5QdeLXW8jnm2PupBs")
}

func main() {
	// addr := "1QHBj5GjAEp7oFKhp5QdeLXW8jnm2PupBs"
	pubkey := "04a4c6b2c8c62ffe17c6740434776c912ff3fc11891922fac4ee93e966817f4446caa82ba4fe6b494f67f870c332fdc5ba22fcdfd6b3a1e88e8bb0693b9bf2c1c2"
	signature := "HKnOcPe/RxF48z5U6JbyetZC7+wmPrlUOumbbecpMVlwbcfGLlTwGBtMDzjD4wxOg/VjQDg7TxHqP/Mfoohp7Cs="
	message := "test"
	verify(pubkey, signature, message)
}
