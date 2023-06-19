## 关于本项目

 Telegram自动开会员源代码 基于 `Golang` ，这是一个完整的24小时全自动开通TG会员的代码，如果你会一点技术，可无缝对接至你的机器人实现24小时自动代开。

## 开始
这只是一个demo版本，没有写配置文件，使用之前，你应该详细阅读源代码里的备注，修改应该代码里备注需要修改的地方。

### 安装环境
项目运行基于`Golang`，你需要先安装`Golang`

+ Windows
  > https://go.dev/ 前往Golang官方网站进行下载安装，如不会建议Google。
+ Linux
  > Centos  
  yum install golang  
  <BR>
  Ubuntu  
  sudo apt-get install golang

### 安装依赖

+ golang
    >   go mod init   
        go mod tidy   
        或者  
        go install

## 实现逻辑

1. 通过解析https://fragment.com的数据
2. 指定开通会员的用户名
3. 解析出相关Payload
4. 携带Payload进行Ton支付
5. 支付完成，开通完成。

## 技术交流/意见反馈

+ MCG技术交流群 https://t.me/MCG_Club

## AD -- 免费领取国际信用卡
>免费领取VISA卡，万事达卡，充值USDT即可随便刷  
可绑微信、支付宝、美区AppStore消费  
24小时自助开卡充值 无需KYC  
无需人工协助，用户可自行免费注册，后台自助实现入金、开卡、绑卡、销卡、查询等操作，支持无限开卡、在线接码。  
✅支持 OpenAi 人工智能 chatGPT PLUS 开通   
✅支持 开通Telegram飞机会员  
➡️➡️➡️ [点击领取你的国际信用卡](https://t.me/EKaPayBot?start=FV6S5XHT9H)

## AD -- 机器人推广

24小时自动发卡机器人：[自动发卡](https://t.me/fakatestbot)
> 24小时自动发卡机器人 对接独角

兑币机 - TRX自动兑换：[兑币机](https://t.me/ConvertTrxBot)
> 自用兑币机，并不是开源版机器人！！！

波场能量机器人：[波场能量机器人](https://t.me/BuyEnergysBot)
> 波场能量租用，有能量时转账USDT不扣TRX，为你节省50-70%的TRX

TG会员秒开机器人：[TG会员秒开-全自动发货](https://t.me/BuySvipBot)
> 24小时自动开通Telegram Premium会员，只需一个用户名即可开通。

## 许可证

根据 MIT 许可证分发。打开 [LICENSE.txt](LICENSE.txt) 查看更多内容。


<p align="right">(<a href="#top">返回顶部</a>)</p>