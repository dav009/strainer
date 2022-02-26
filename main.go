package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/dav009/strainer/ergo"
	"github.com/jwalton/gchalk"
)

func fetchLatestTx(url string, lastHeight float32) (float32, error) {
	node := ergo.Node{Url: url}
	latestHeight, err := node.LastHeight()
	if err != nil {
		return 0, err
	}
	if latestHeight > lastHeight {

		headerIds, err := node.MainChainHeaderIdAtHeight(latestHeight)
		if err != nil {
			return 0, err
		}
		fmt.Println(headerIds)
		block, err := node.TxsAtHeader(headerIds[0])
		if err != nil {
			return 0, err
		}

		for _, tx := range block.BlockTransactions.Transactions {
			printTx(tx)
		}
		return latestHeight, nil
	}
	return lastHeight, nil
}

func listen(url string, done chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		latestHeight := float32(-1)
		for {
			select {
			case <-done:
				fmt.Println("closing ticket..")
				return
			case <-ticker.C:
				//fmt.Println("Tick at", t)
				latestHeight2, err := fetchLatestTx(url, latestHeight)
				if err != nil {
					panic(err)
				}
				latestHeight = latestHeight2
			}
		}
	}()

}

func printTx(tx ergo.Transaction) {
	prefix := fmt.Sprintf("Block:XXX>TX: %s", tx.Id[:8])
	orange := gchalk.RGB(255, 136, 0)

	fmt.Printf("%s\t%s\t%s\t%s\n",
		gchalk.Black(prefix),
		gchalk.Blue(fmt.Sprintf("■ TX")),
		gchalk.White(reflect.TypeOf(tx).Name()),
		gchalk.White(fmt.Sprintf("%s %d", tx.Id, tx.Size)),
	)

	for _, input := range tx.Inputs {
		metadataType := "UTXO"
		subType := reflect.TypeOf(input).Name()
		ResourceTypePrefix := fmt.Sprintf("■ %s", metadataType)
		fmt.Printf("%s\t%s\t%s\t%s\n",
			gchalk.Black(prefix),
			orange(ResourceTypePrefix),
			gchalk.White(subType),
			gchalk.White(fmt.Sprintf("%+v", input)),
		)
	}

	for _, input := range tx.Inputs {
		metadataType := "UTXO"
		subType := reflect.TypeOf(input).Name()
		ResourceTypePrefix := fmt.Sprintf("■ %s", metadataType)
		fmt.Printf("%s\t%s\t%s\t%s\n",
			gchalk.Black(prefix),
			orange(ResourceTypePrefix),
			gchalk.White(subType),
			gchalk.White(fmt.Sprintf("%+v", input)),
		)
	}

	for _, output := range tx.Outputs {
		metadataType := "UTXO"
		subType := reflect.TypeOf(output).Name()
		ResourceTypePrefix := fmt.Sprintf("■ %s", metadataType)
		fmt.Printf("%s\t%s\t%s\t%s\n",
			gchalk.Black(prefix),
			orange(ResourceTypePrefix),
			gchalk.White(subType),
			gchalk.White(fmt.Sprintf("%+v", output)),
		)
	}

}

func main() {
	url := os.Args[1]
	done := make(chan bool, 1)
	var quit = make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		done <- true
		close(quit)
		fmt.Printf("You pressed ctrl + C. User interrupted infinite loop.")
		os.Exit(0)
	}()
	listen(url, done)
	<-quit
}
