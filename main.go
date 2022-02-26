package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Node struct {
	url string
}

type PowSolution struct {
	Pk string
	W  string
	N  string
	D  float64
}

type Input struct {
	BoxId         string
	SpendingProof string `json:"-"`
}

type Output struct {
	BoxId               string
	value               int32
	assets              []string `json:"-"`
	creationHeight      int32
	additionalRegisters string `json: "-"`
	transactionId       string
	index               int16
}

type Transaction struct {
	Id         string
	Inputs     []Input
	DataInputs []string `json:"-"`
	Outputs    []Output
	Size       int16
}

type BlockTransactions struct {
	HeaderId     string
	Transactions []Transaction
	BlockVersion int16
	size         int16
}

type Header struct {
	ExtensionId      string
	Difficulty       string
	Votes            string
	Timestamp        float32
	Size             int
	StateRoot        string
	Height           float32
	NBits            float32
	Version          int16
	Id               string
	AdProofsRoot     string
	TransactionsRoot string
	ExtensionHash    string
	PowSolutions     PowSolution
	AdProofsId       string
	TransactionsId   string
	ParentId         string
}

type Block struct {
	Header            Header
	BlockTransactions BlockTransactions
	Extension         string `json:"-"`
	AdProofs          string `json:"-"`
	size              int16
}

type Info struct {
	FullHeight float32
}

func (n *Node) txsAtHeader(headerId string) (Block, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", n.url, "blocks", headerId)
	fmt.Println(endpoint)
	var block Block
	r, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return block, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Charset", "utf-8")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return block, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyBytes, &block)

	if err != nil {
		return block, err
	}
	return block, err
}

func (n *Node) mainChainHeaderIdAtHeight(height float32) ([]string, error) {
	endpoint := fmt.Sprintf("%s/%s/%f", n.url, "blocks/at", height)
	fmt.Println(endpoint)
	r, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Charset", "utf-8")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	var headers []string
	err = json.Unmarshal(bodyBytes, &headers)
	if err != nil {
		return nil, err
	}
	if len(headers) != 1 {
		return nil, errors.New(fmt.Sprintf("error: no headers at height: %d ", height))
	}
	return headers, err

}

func (n *Node) lastHeight() (float32, error) {
	endpoint := fmt.Sprintf("%s/info", n.url)
	fmt.Println(endpoint)
	r, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return 0, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Charset", "utf-8")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	var info Info
	err = json.Unmarshal(bodyBytes, &info)
	if err != nil {
		return 0, err
	}

	return info.FullHeight, nil
}

func fetchLatestTx(lastHeight float32) (float32, error) {
	node := Node{url: "http://192.168.3.99:9053"}
	latestHeight, err := node.lastHeight()
	if err != nil {
		return 0, err
	}
	if latestHeight > lastHeight {

		headerIds, err := node.mainChainHeaderIdAtHeight(latestHeight)
		if err != nil {
			return 0, err
		}
		fmt.Println(headerIds)
		block, err := node.txsAtHeader(headerIds[0])
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

func listen(done chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		latestHeight := float32(-1)
		for {
			select {
			case <-done:
				fmt.Println("closing ticket..")
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				latestHeight2, err := fetchLatestTx(latestHeight)
				if err != nil {
					panic(err)
				}
				latestHeight = latestHeight2
			}
		}
	}()

}

func printTx(tx Transaction) {
	fmt.Println(tx.Id)
}

func main() {

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
	listen(done)
	<-quit
}
