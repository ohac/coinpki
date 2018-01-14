package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/itchyny/base58-go"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

func sha256d(body []byte) []byte {
	msghash1 := sha256.Sum256([]byte(body))
	msghash2 := sha256.Sum256(msghash1[:])
	return msghash2[:]
}

const (
	H_BTC  string = "Bitcoin Signed Message:\n"
	H_MONA string = "Monacoin Signed Message:\n"
)

type CoinParams struct {
	Header string
	Magic  int
}

func messagehash(message, header string) (msghash2 []byte, err error) {
	hlen := len(header)
	if hlen >= 0xfd {
		err = fmt.Errorf("long header is not supported")
		return
	}
	mlen := len(message)
	if mlen >= 0xfd {
		err = fmt.Errorf("long message is not supported")
		return
	}
	btcmsg := string([]byte{byte(hlen)})
	btcmsg += header
	btcmsg += string([]byte{byte(mlen)})
	btcmsg += message
	msghash2 = sha256d([]byte(btcmsg))
	return
}

func parse_signature(signature string) (sig Signature, recid int, err error) {
	sigraw, err2 := base64.StdEncoding.DecodeString(signature)
	if err2 != nil {
		err = err2
		return
	}
	r0 := sigraw[0] - 27
	recid = int(r0 & 3)
	compressed := (r0 & 4) == 1
	if compressed {
		err = fmt.Errorf("compressed type is not supported")
		return
	}
	sig.R.SetBytes(sigraw[1 : 1+32])
	sig.S.SetBytes(sigraw[1+32 : 1+32+32])
	return
}

func pubtoaddr(pubkey_xy2 XY, compressed bool, magic int) (bcpy []byte) {
	size := 65
	if compressed {
		size = 33
	}
	out := make([]byte, size)
	pubkey_xy2.GetPublicKey(out)
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(out)
	pub_hash_1 := sha256_h.Sum(nil)
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)
	bcpy = append([]byte{byte(magic)}, pub_hash_2...)
	hash2 := sha256d(bcpy)
	bcpy = append(bcpy, hash2[0:4]...)
	return
}

func addrtostr(bcpy []byte) (s string, err error) {
	z := new(big.Int)
	z.SetBytes(bcpy)
	enc := base58.BitcoinEncoding
	var encdd []byte
	encdd, err = enc.Encode([]byte(z.String()))
	if err != nil {
		return
	}
	s = string(encdd)
	for _, v := range bcpy {
		if v != 0 {
			break
		}
		s = "1" + s
	}
	return
}

func sigmestoaddr(signature, message string, params CoinParams,
	compressed bool) (addr string, err error) {
	msghash2, err2 := messagehash(message, params.Header)
	if err2 != nil {
		err = err2
		return
	}
	sig, recid, err2 := parse_signature(signature)
	if err2 != nil {
		err = err2
		return
	}
	var msg Number
	msg.SetBytes(msghash2)

	var pubkey_xy2 XY
	ret2 := sig.recover(&pubkey_xy2, &msg, recid)
	if !ret2 {
		err = fmt.Errorf("recover pubkey failed")
		return
	}

	bcpy := pubtoaddr(pubkey_xy2, compressed, params.Magic)
	s, err2 := addrtostr(bcpy)
	if err2 != nil {
		err = err2
		return
	}
	addr = s
	return
}

func verify(addr1, signature, message string, params CoinParams,
	compressed bool) {
	addr2, err := sigmestoaddr(signature, message, params, compressed)
	if err != nil {
		println("failed")
		return
	}
	if addr1 == addr2 {
		println("verified")
	} else {
		println("failed", addr1, addr2)
	}
}

func main() {
	P_BTC := CoinParams{
		Header: H_BTC,
		Magic:  0}
	P_MONA := CoinParams{
		Header: H_MONA,
		Magic:  50}
	verify("1QHBj5GjAEp7oFKhp5QdeLXW8jnm2PupBs",
		"HKnOcPe/RxF48z5U6JbyetZC7+wmPrlUOumbbecpMVlwbcfGLlTwGBtMDzjD4wxOg/VjQDg7TxHqP/Mfoohp7Cs=",
		"test", P_BTC, false)
	verify("MQ8q9jSGQdnHmZe4kfjGUkzHoPF9GCBbN6",
		"IKO8h8iYp0wIBCh+D+/ixJ2MovYueUZDsFuvcvIPqNFnGTtL/eggy7HNymCbKemHbLR0QB1DpC6o6/By/eubXzI=",
		"test", P_MONA, true)
}
