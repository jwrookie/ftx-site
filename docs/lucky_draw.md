# lucky draw


### CSRF-TOKEN
> 所有的接口，在请求头中必须加上CSRF-TOKEN。
> 
> 如何获取？
> 
> 按公钥对当前时间戳进行RSA256，公钥文件：ftx-site/config/rsa_public_key.pem
> 后端会解析出时间戳，并判断当前时间 与 请求中携带的时间是否在一个合理范围内，如果是，则认为这是一个正常的请求
> 
> 注意：`公钥文件` 与 `加密逻辑`，前端应该使用 不可逆混淆 来防止代码泄露


## 获取抽奖资格
```curl
curl --location --request POST 'http://127.0.0.1:8080/lucky/token' \
--header 'CSRF-TOKEN: lwhbjvf4hiqgeulbakjrq54fwelfn11ksdfj65ksdg63lgrndlkKE2FJLFK' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "1231@gmail.com",
    "kyc_level": "KYC0",
    "personality":"IATC"
}' | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   425  100   341  100    84  49199  12119 --:--:-- --:--:-- --:--:--  415k
{
  "code": 0,
  "msg": "ok",
  "data": {
    "token": "cXHaUYoIJayHJ52wFc0M1t1MyD17NBpz6S7073Tren4CnZ00dv91w+Em+ZRG/TD0Gh5gvU9GnRcrkqaxGgFKtxR9epKfhV0puisT9CWBVS+tTm4wzvnCGEvPNkwaYNQqHmicXzPMOOz/hQVlIJx41QAVQ4aXsjfYzn8holTHIQK8wwo6qPT0o0F6M8pEQ8gHWXZdKZcYkaHizYVV/PUTAdcCccCFIjxDf8pxFClg8t0XKOv7O0SLriFrlnfaPdDAo1FMEhQhie+L41Hedo+/yQlqM4Vdskm0MOYyHPXnSRI="
  }
}
```
> payload：
> 
> email: string, required，email format
> 
> kyc_level: string, required，one of `KYC0` `KYC1` `KYC2`
> 
> personality: string, required，one of `IATC` `EATC` `IATM` `EATM` `IAFC` `EAFC` `IAFM` `EAFM` `IPTC` `EPTC` `IPTM` `EPTM` `IPFC` `EPFC` `IPFM` `EPTM`


## 进行抽奖
> header 中的`DRAW-TOKEN` 是获取抽奖资格时拿到的`token`
```curl
curl --location --request POST 'http://127.0.0.1:8080/lucky/draw' \
--header 'DRAW-TOKEN: cXHaUYoIJayHJ52wFc0M1t1MyD17NBpz6S7073Tren4CnZ00dv91w+Em+ZRG/TD0Gh5gvU9GnRcrkqaxGgFKtxR9epKfhV0puisT9CWBVS+tTm4wzvnCGEvPNkwaYNQqHmicXzPMOOz/hQVlIJx41QAVQ4aXsjfYzn8holTHIQK8wwo6qPT0o0F6M8pEQ8gHWXZdKZcYkaHizYVV/PUTAdcCccCFIjxDf8pxFClg8t0XKOv7O0SLriFrlnfaPdDAo1FMEhQhie+L41Hedo+/yQlqM4Vdskm0MOYyHPXnSRI=' \
--header 'CSRF-TOKEN: lwhbjvf4hiqgeulbakjrq54fwelfn11ksdfj65ksdg63lgrndlkKE2FJLFK' | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    53  100    53    0     0   1060      0 --:--:-- --:--:-- --:--:--  1232
{
  "code": 0,
  "msg": "ok",
  "data": {
    "prize": "FTX束口袋背包"
  }
}
```


## 领取奖品
> header 中的`DRAW-TOKEN` 是获取抽奖资格时拿到的`token`
```curl
curl --location --request POST 'http://127.0.0.1:8080/lucky/award' \
--header 'DRAW-TOKEN: cXHaUYoIJayHJ52wFc0M1t1MyD17NBpz6S7073Tren4CnZ00dv91w+Em+ZRG/TD0Gh5gvU9GnRcrkqaxGgFKtxR9epKfhV0puisT9CWBVS+tTm4wzvnCGEvPNkwaYNQqHmicXzPMOOz/hQVlIJx41QAVQ4aXsjfYzn8holTHIQK8wwo6qPT0o0F6M8pEQ8gHWXZdKZcYkaHizYVV/PUTAdcCccCFIjxDf8pxFClg8t0XKOv7O0SLriFrlnfaPdDAo1FMEhQhie+L41Hedo+/yQlqM4Vdskm0MOYyHPXnSRI=' \
--header 'CSRF-TOKEN: lwhbjvf4hiqgeulbakjrq54fwelfn11ksdfj65ksdg63lgrndlkKE2FJLFK' \
--header 'Content-Type: application/json' \
--data-raw '{
    "prize": "FTX束口袋背包",
    "user_name": "jw",
    "user_phone": "12311112222",
    "address": "xxx"
}' | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   148  100    33  100   115    787   2745 --:--:-- --:--:-- --:--:--  4228
{
  "code": 0,
  "msg": "ok",
  "data": null
}
```
> payload：
> 
> prize: string, required, one of `FTX x MLB棒球外套` `FTX束口袋背包` `FTX棒球帽` `FTX灰色T恤` `交易手续费抵扣券10USD` `FTX小龙人暖手充电宝` `FTX雪花真空杯` `FTX超萌小耳朵发箍 + 小金勺子`
> 
> clothes_size: string, required when prize is FTX灰色T恤
> 
> user_name: string, required
> 
> user_phone: string, required
> 
> address: string, required


### 获取奖池金额
```curl
curl --location --request GET 'http://127.0.0.1:8080/lucky/jackpot' \
--header 'CSRF-TOKEN: lwhbjvf4hiqgeulbakjrq54fwelfn11ksdfj65ksdg63lgrndlkKE2FJLFK' | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    47  100    47    0     0   7170      0 --:--:-- --:--:-- --:--:-- 47000
{
  "code": 0,
  "msg": "ok",
  "data": {
    "jackpot": "5003"
  }
}
```


## 查询中奖信息
```curl
curl --location --request GET 'http://127.0.0.1:8080/lucky/1231@gmail.com' \
--header 'CSRF-TOKEN: lwhbjvf4hiqgeulbakjrq54fwelfn11ksdfj65ksdg63lgrndlkKE2FJLFK' | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   303  100   303    0     0  21676      0 --:--:-- --:--:-- --:--:-- 37875
{
  "code": 0,
  "msg": "ok",
  "data": {
    "lucky_id": "417058349658414957",
    "email": "1231@gmail.com",
    "kyc_level": "KYC0",
    "personality": "IATC",
    "prize": "FTX束口袋背包",
    "clothes_size": "",
    "user_name": "jw",
    "user_phone": "12311112222",
    "address": "xxx",
    "created_at": 1658115747822,
    "updated_at": 1658115747822,
    "deleted_at": 0
  }
}
```