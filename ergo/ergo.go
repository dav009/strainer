package ergo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Node struct {
	Url string
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
	Value               int64
	Assets              []string `json:"-"`
	CreationHeight      int32
	additionalRegisters string `json: "-"`
	TransactionId       string
	Index               int16
	ErgoTree            string
}

type Transaction struct {
	Id         string
	Inputs     []Input
	DataInputs []string `json:"-"`
	Outputs    []Output
	Size       int32
}

type BlockTransactions struct {
	HeaderId     string
	Transactions []Transaction
	BlockVersion int32
	size         int32
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

func (n *Node) TxsAtHeader(headerId string) (Block, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", n.Url, "blocks", headerId)
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

func (n *Node) MainChainHeaderIdAtHeight(height float32) ([]string, error) {
	endpoint := fmt.Sprintf("%s/%s/%f", n.Url, "blocks/at", height)
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

func (n *Node) LastHeight() (float32, error) {
	endpoint := fmt.Sprintf("%s/info", n.Url)
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
