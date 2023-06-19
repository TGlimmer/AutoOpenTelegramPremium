package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type PaymentService struct {
}

type RequestPayload struct {
	Query      string
	Months     int
	Recipient  string
	ID         string
	ShowSender int
	Method     string
}

const apiURL = "https://fragment.com/api?hash=替换为你自己的hash" //替换为你自己的hash

func main() {

	// 创建 PaymentService 实例
	ps := &PaymentService{}

	// 示例 1：获取对方信息
	payload1 := RequestPayload{
		Query:  "miya0v0", //需要开通的用户名
		Months: 3,         //需要开通的月份
		Method: "searchPremiumGiftRecipient",
	}
	result1, err := ps.SendRequest(payload1)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	recipient := result1["recipient"].(string)
	log.Println("Result recipient:", recipient)
	// 获取 reqid
	payload2 := RequestPayload{
		Recipient: recipient,
		Months:    3, //需要开通的月份
		Method:    "initGiftPremiumRequest",
	}
	result2, err := ps.SendRequest(payload2)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	reqId := result2["req_id"].(string)
	log.Println("Result reqId:", reqId)
	// 确认订单
	ConfirmOrder := RequestPayload{
		ID:         reqId,
		ShowSender: 1,
		Method:     "getGiftPremiumLink",
	}
	ConfirmOrderResult, err := ps.SendRequest(ConfirmOrder)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	currentTime := time.Now()
	expireTime := ConfirmOrderResult["expire_after"].(float64)
	fiveMinutesLater := currentTime.Add(time.Duration(expireTime) * time.Second)
	log.Println("订单过期时间:", fiveMinutesLater.Format("2006-01-02 15:04:05"))
	// 获取payload
	payload, amount, err := ps.GetRawRequest(reqId)
	if err != nil {
		log.Println("Error:", err)
	}
	log.Println("Result payload:", payload)
	transferTon(amount, payload)

}

func (ps *PaymentService) SendRequest(payload RequestPayload) (map[string]interface{}, error) {
	client := &http.Client{}

	data := url.Values{}
	if payload.Query != "" {
		data.Set("query", payload.Query)
	}
	if payload.Months != 0 {
		data.Set("months", strconv.Itoa(payload.Months))
	}
	if payload.Recipient != "" {
		data.Set("recipient", payload.Recipient)
	}
	if payload.ID != "" {
		data.Set("id", payload.ID)
	}
	if payload.ShowSender != 0 {
		data.Set("show_sender", strconv.Itoa(payload.ShowSender))
	}
	data.Set("method", payload.Method)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "这里需要替换为你自己的cookie") //需要替换为你自己的 https://fragment.com cookie

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	method := data.Get("method")
	switch method {
	case "searchPremiumGiftRecipient":
		if ok, exists := result["ok"]; !exists || !ok.(bool) {
			return nil, errors.New("请求失败")
		}
		found, ok := result["found"].(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid 'found' field in the response")
		}
		recipient, ok := found["recipient"].(string)
		if !ok {
			return nil, errors.New("invalid 'recipient' field in the 'found' object")
		}
		log.Printf("recipient: %v\n", recipient)
		return found, nil
	case "initGiftPremiumRequest":
		req_id, ok := result["req_id"].(string)
		if !ok {
			return nil, errors.New("invalid 'req_id' field in the 'found' object")
		}
		log.Printf("req_id: %v\n", req_id)
		return result, nil
	case "getGiftPremiumLink":
		if ok, exists := result["ok"]; !exists || !ok.(bool) {
			return nil, errors.New("请求失败")
		}
		//fmt.Printf("result: %v\n", result)
		_, ok := result["expire_after"].(float64)
		if !ok {
			return nil, errors.New("invalid 'expire_after' field in the 'found' object")
		}
		return result, nil
	default:
		return nil, errors.New("invalid method")
	}

}

func (p *PaymentService) GetRawRequest(id string) (string, string, error) {
	url := fmt.Sprintf("https://fragment.com/tonkeeper/rawRequest?id=%s", id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "这里需要替换为你自己的cookie") //需要替换为你自己的 https://fragment.com cookie
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	bodyObj, ok := result["body"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("invalid 'body' field in the response")
	}

	params, ok := bodyObj["params"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("invalid 'params' field in the 'body' object")
	}

	messages, ok := params["messages"].([]interface{})
	if !ok || len(messages) == 0 {
		return "", "", errors.New("invalid 'messages' field in the 'params' object")
	}

	firstMessage, ok := messages[0].(map[string]interface{})
	if !ok {
		return "", "", errors.New("invalid 'messages' field in the 'params' object")
	}
	amount, ok := firstMessage["amount"].(float64)
	payamount := fmt.Sprintf("%g", amount/1e9)
	if !ok {

		return "", "", errors.New("invalid 'amount' field in the first message object")
	}
	payload, ok := firstMessage["payload"].(string)
	if !ok {
		return "", "", errors.New("invalid 'payload' field in the first message object")
	}

	log.Printf("支付金额: %v\n", payamount)
	// base64 decode payload
	padding := len(payload) % 4
	if padding > 0 {
		payload += strings.Repeat("=", 4-padding)
	}
	log.Printf("payload: %v\n", payload)
	decodedPayload, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		fmt.Println("Error:", err)
	}
	decodedPayloadStr := string(decodedPayload)
	// 使用正则表达式匹配目标字符串
	match := extractRefFromPayload(decodedPayloadStr)
	if len(match) == 8 {
		log.Println("提取到的字符:", match)
	} else {
		log.Println("未找到匹配的字符")
	}
	return match, payamount, nil
}

func extractRefFromPayload(payload string) string {
	decodedPayloadStr := payload

	refStr := ""
	i := strings.Index(decodedPayloadStr, "#")

	//fmt.Println("Input:", decodedPayloadStr)

	if i != -1 {
		// 开始从 "Ref#" 之后的位置提取字符
		for i = i + 1; i < len(decodedPayloadStr) && len(refStr) < 8; i++ {
			if unicode.IsLetter(rune(decodedPayloadStr[i])) || unicode.IsNumber(rune(decodedPayloadStr[i])) {
				refStr += string(decodedPayloadStr[i])
			}
		}
	} else {
		fmt.Println("Ref# not found")
	}

	return refStr
}

func transferTon(amount string, payload string) {
	client := liteclient.NewConnectionPool()

	// 这里使用官方的节点配置文件 比较慢 你可以自创节点
	configUrl := "https://ton-blockchain.github.io/global.config.json"
	err := client.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		panic(err)
	}

	// 初始化 ton api lite 连接包装器
	api := ton.NewAPIClient(client)

	// 如果我们想将所有请求路由到同一个节点，我们可以使用它
	ctx := client.StickyContext(context.Background())

	// 我们需要新的区块信息来运行 get 方法
	b, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalln("get block err:", err.Error())
		return
	}

	addr := address.MustParseAddr("EQBAjaOyi2wGWlk-EDkSabqqnF-MrrwMadnwqrurKpkla9nE") //这里是官方的钱包地址 请勿更改 否则开通会员将不会到账

	// 我们使用 WaitForBlock 来确保块准备好，
	// 它是可选的，但可以让我们避免 liteserver 块未准备好错误
	res, err := api.WaitForBlock(b.SeqNo).GetAccount(ctx, b, addr)
	if err != nil {
		log.Fatalln("get account err:", err.Error())
		return
	}

	log.Printf("Is active: %v\n", res.IsActive)
	if res.IsActive {
		log.Printf("Status: %s\n", res.State.Status)
		log.Printf("Balance: %s TON\n", res.State.Balance.TON())
		if res.Data != nil {
			log.Printf("Data: %s\n", res.Data.Dump())
		}
	}

	fmt.Printf("\nTransactions:\n")

	// 这里填写你自己的TON钱包的助记词 请用空格分割
	words := strings.Split("这里填写你自己的TON钱包的助记词", " ")

	w, err := wallet.FromSeed(api, words, wallet.V4R2)

	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("wallet address:", w.Address())

	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		log.Fatalln("CurrentMasterchainInfo err:", err.Error())
		return
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	if balance.NanoTON().Uint64() >= 3000000 {
		addr := address.MustParseAddr("EQBAjaOyi2wGWlk-EDkSabqqnF-MrrwMadnwqrurKpkla9nE") //这里是官方的钱包地址 请勿更改 否则开通会员将不会到账
		// 创建消息
		comment, err := wallet.CreateCommentCell(fmt.Sprintf("Telegram Premium for 3 months Ref#%s", payload)) //这里的3是开通3个月需要和上方的月份对应
		if err != nil {
			log.Fatalln("CreateComment err:", err.Error())
			return
		}
		amount := tlb.MustFromTON(amount)

		log.Printf("付款地址: %s , 金额: %v\n", addr, amount)

		var messages []*wallet.Message
		messages = append(messages, &wallet.Message{
			Mode: 1,
			InternalMessage: &tlb.InternalMessage{
				Bounce:  false, // 不允许回退 一般情况下不需要回退
				DstAddr: addr,
				Amount:  amount,
				Body:    comment,
			},
		})

		log.Println("发送付款请求并等待区块确认中...")

		// 发送交易 并等待区块确认
		tx, block, err := w.SendManyWaitTransaction(context.Background(), messages)
		if err != nil {
			log.Fatalln("Transfer err:", err.Error())
			return
		}
		balance, err = w.GetBalance(context.Background(), block)
		if err != nil {
			log.Fatalln("GetBalance err:", err.Error())
			return
		}

		log.Println("区块已确认,交易成功! 交易hash:", base64.StdEncoding.EncodeToString(tx.Hash))
		log.Println("查看交易: https://tonscan.org/tx/" + base64.URLEncoding.EncodeToString(tx.Hash))
		log.Println("钱包还剩余余额: ", balance.TON())
		return
	}

	log.Println("not enough balance:", balance.TON())
}
