package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/hectorchu/gonano/rpc"
	"github.com/hectorchu/gonano/util"
	"github.com/hectorchu/gonano/wallet"
)

func main() {
	var (
		readAddress = flag.String("address", "", "read file from NANO address")
		writeFile   = flag.String("file", "", "write file to a NANO address")
		rpcURL      = flag.String("rpc", "https://mynano.ninja/api/node", "RPC URL to use")
		err         = errors.New("Specify either -address or -file")
	)
	flag.Parse()
	switch {
	case *readAddress != "" && *writeFile != "":
	case *readAddress != "":
		err = read(*readAddress, *rpcURL)
	case *writeFile != "":
		err = write(*writeFile, *rpcURL)
	default:
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func read(address, rpcURL string) (err error) {
	client := rpc.Client{URL: rpcURL}
	info, err := client.AccountInfo(address)
	if err != nil {
		return
	}
	hashes, err := client.Chain(info.Frontier, -1)
	if err != nil {
		return
	}
	for i := len(hashes) - 2; i >= 0; i-- {
		info, err := client.BlockInfo(hashes[i])
		if err != nil {
			return err
		}
		data, err := util.AddressToPubkey(info.Contents.Representative)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	}
	return
}

func write(name, rpcURL string) (err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return
	}
	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		return
	}
	w, err := wallet.NewWallet(buf)
	if err != nil {
		return
	}
	w.RPC.URL = rpcURL
	a, err := w.NewAccount(nil)
	if err != nil {
		return
	}
	fmt.Println("Deposit some NANO to", a.Address())
	for {
		balance, pending, err := a.Balance()
		if err != nil {
			return err
		}
		var zero big.Int
		if balance.Cmp(&zero) > 0 || pending.Cmp(&zero) > 0 {
			break
		}
		time.Sleep(5 * time.Second)
	}
	err = a.ReceivePendings()
	if err != nil {
		return
	}
	bar := pb.StartNew(int(fi.Size()))
	for {
		buf := make([]byte, 32)
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		address, err := util.PubkeyToAddress(buf)
		if err != nil {
			return err
		}
		_, err = a.ChangeRep(address)
		if err != nil {
			return err
		}
		bar.Add(n)
	}
	bar.Finish()
	fmt.Println("Stored", f.Name(), "to", a.Address())
	return
}
