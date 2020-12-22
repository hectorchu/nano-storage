package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/hectorchu/gonano/util"
	"github.com/hectorchu/gonano/wallet"
)

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: nano-storage <file to store>")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		fatal(err)
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		fatal(err)
	}
	w, err := wallet.NewWallet(buf)
	if err != nil {
		fatal(err)
	}
	a, err := w.NewAccount(nil)
	if err != nil {
		fatal(err)
	}
	fmt.Println("Deposit some NANO to", a.Address())
	for {
		balance, pending, err := a.Balance()
		if err != nil {
			fatal(err)
		}
		var zero big.Int
		if balance.Cmp(&zero) > 0 || pending.Cmp(&zero) > 0 {
			break
		}
		time.Sleep(5 * time.Second)
	}
	err = a.ReceivePendings()
	if err != nil {
		fatal(err)
	}
	bar := pb.StartNew(int(fi.Size()))
	for {
		buf := make([]byte, 32)
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fatal(err)
		}
		address, err := util.PubkeyToAddress(buf)
		if err != nil {
			fatal(err)
		}
		_, err = a.ChangeRep(address)
		if err != nil {
			fatal(err)
		}
		bar.Add(n)
	}
	bar.Finish()
	fmt.Println("Stored", f.Name(), "to", a.Address())
}
