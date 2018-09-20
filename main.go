package main

import (
	"fmt"
	"github.com/vsergeev/btckeygenie/btckey"
	"crypto/rand"
	"log"
	"net/http"
	"strconv"
	"io/ioutil"
	"raiders/model"
	"encoding/json"
	"time"
)

var (
	totalReceived int
	priv btckey.PrivateKey
	err error
	addressCompressed string
)

func main() {
	count := 0
	for {
		t := time.Now()
		// private key 생성
		priv, err = btckey.GenerateKey(rand.Reader)
		if err != nil {
			log.Fatalf("Generating keypair: %s\n", err)
		}
		addressCompressed = priv.ToAddress()

		fmt.Println("--------------------------------------------------------------------------------------")
		fmt.Println("[" + strconv.Itoa(count)+ "] " + t.String())
		fmt.Println("[" + strconv.Itoa(count)+ "] " + "Private Key: " + priv.ToWIF())

		// 비트코인을 찾은 경우 break 테스트
		//if count == 3 {
		//	addressCompressed = "1JdTWTAubsDWXWd7wcsFwMuSapMBh6efrQ"
		//}

		fmt.Println("[" + strconv.Itoa(count)+ "] " + "address: " + addressCompressed)
		totalReceived = getTotalReceived(addressCompressed, count)
		if totalReceived > 0 {
			fmt.Println("I found Bitcoin!")
			break
		}
		count++;
	}
	fmt.Println("Ends the search.")
}

func getTotalReceived(address string, count int) int {
	// 조회 가능 api 3군데
	// https://insight.bitpay.com/api/addr/1Nq8PMyYKmtED92zBB4moNARjUmVaShwcH?noTxList=1
	// https://blockexplorer.com/api/addr/1Nq8PMyYKmtED92zBB4moNARjUmVaShwcH?noTxList=1
	// https://blockchain.info/rawaddr/1Nq8PMyYKmtED92zBB4moNARjUmVaShwcH
	// https://api.blockcypher.com/v1/btc/main/addrs/1DEP8i3QJCsomS4BSMY2RpU1upv62aGvhD/balance
	var url string

	switch count % 4 {
		case 0:
			// blockchain.info
			//fmt.Println("blockchain.info 에서 조회중...")
			url = "https://blockchain.info/rawaddr/" + address
			break
		case 1:
			// insight api explorer 인 경우
			//fmt.Println("blockexplorer 에서 조회중...")
			url = "https://blockexplorer.com/api/addr/" + address + "?noTxList=1"
			break
		case 2:
			// insight api explorer 인 경우
			//fmt.Println("bitpay 에서 조회중...")
			url = "https://insight.bitpay.com/api/addr/" + address + "?noTxList=1"
			break
		case 3:
			// blockcypher api 인 경우
			//fmt.Println("blockcypher 에서 조회중...")
			url = "https://api.blockcypher.com/v1/btc/main/addrs/" + address + "/balance"
			break
	}

	// url := "https://blockchain.info/rawaddr/" + address
	req, _ := http.NewRequest("GET", url, nil)
	res, error := http.DefaultClient.Do(req)
	if error != nil {
		fmt.Println("Balance API Error: " + strconv.Itoa(count))
		return -1
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var result int

	switch count % 4 {
	case 0:
		// blockchain.info 인 경우
		var addrJsonStruct model.AddrJsonStruct
		_ = json.Unmarshal(body, &addrJsonStruct)
		result = addrJsonStruct.TotalReceived
		break
	case 1:
	case 2:
		// insight api explorer 인 경우
		var addrJsonStruct model.AddrJsonStructByInsight
		_ = json.Unmarshal(body, &addrJsonStruct)
		result = addrJsonStruct.TotalReceivedSat
		break
	case 3:
		// blockcypher 인 경우
		var addrJsonStruct model.AutoGeneratedByBlockcypher
		_ = json.Unmarshal(body, &addrJsonStruct)
		result = addrJsonStruct.TotalReceived
		break
	}
	fmt.Println("TotalReceived: " + strconv.Itoa(result))

	return result
}