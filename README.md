# GoAuth

## Overview

該專案包含兩個核心功能，帳號系統及商品推薦系統，帳號系統提供帳號註冊、登入、驗證等功能，商品推薦系統則實作資料庫存取、快取、單次請求等效能優化技巧

## 設計說明

### Authentication

- 在 login 後回傳給 client 的 token，使用的是 `JWT Token`，並且在 token 中加入了 access token 的過期時間，在 middleware 中檢查 JWT Token 有效性及過期時間
- ACCESS_TOKEN_EXP_MINUTES 是作為環境變數存在，方便在不同環境下可以有不同時長的 token，例如在開發環境下可以設定較長的時間，減少替換成本，而在營運環境下設定較短時間，來提升安全性
- 在註冊的信箱驗證碼上，將 uid 及 unix timestamp 資料透過 `AES` 對稱式加密，並且在解密時檢查是否為正確的驗證碼，後端伺服器不需要儲存驗證碼，減少資料庫負擔及維護成本

### Account

- 在帳號系統中，使用了 `bcrypt` 來對密碼進行 hash 保護使用者的密碼
- 目前是以信箱作為帳號，但考慮到未來可能會有其他帳號方式，所以在資料庫設計上使用了 uid 來作為帳號的唯一識別碼
- 在帳號系統中，使用了軟刪除 `deleted_at` 欄位來標記帳號是否被刪除，而不是直接刪除資料，這樣可以保留刪除帳號的紀錄，並且在未來可能會有復原帳號的需求

### Cache

- 在熱門資料存取上，採用了 `Cache Aside` 模式來提升效能，當快取失效時，會向資料庫取得資料，並且在取得資料後，將資料存入快取中，以提升效能
- 為了減緩當快取資料過期的期間，請求會重複的向資料庫取得資料，使用 `single flight` 來避免重複的資料庫存取，提升效能及減少資源浪費

### Clean Architecture

- 參考了 [Clean Architecture](https://github.com/bxcodec/go-clean-arch) 來設計，將程式碼分為不同層級，並且將依賴性從外部注入，以達到程式碼可測試、可維護、可擴展等目的
- 原設計將每個業務場景都獨立出各自的 repository、delivery、usecase 等層級，但在實作時發現，較少有 repository、delivery、usecase 獨自測試的需求，所以這部分並未採用，僅參考 domain 設計、依賴相依方向由外而內，以減少程式碼複雜度

- 以 Send email service 為例，在 main.go 的初始化時，將 Send email service 的實作注入到服務中，以達到`Dependency Injection(DI)`的目的。可以在不同情境決定要使用螢幕輸出、寄信服務、或是其他方式來實作 Send email service

### Custom Error

將錯誤訊息統一化，並且提供給 client 一致的錯誤訊息格式，方便 client 進行錯誤處理，並且在後端程式碼中，可以更容易的進行錯誤處理

## Project Structure

### `config`

## 所需軟體/服務

- Go 1.22
- MySQL 8.0
- Redis

### 環境設定

1. 開啟 `config/local.sh`
2. 修改資料庫連線方式等

### 執行

```shell
make run
```

### 使用

- 註冊帳號

```shell
curl 'localhost:9030/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "go@com.com",
    "password": "Password1~"
}'
```

- 驗證帳號

```shell
curl --location 'localhost:9030/verify-email' \
--header 'Content-Type: application/json' \
--data '{
    "verification_code": "從螢幕上取得驗證碼"
}'
```

- 登入

```shell
curl --location 'localhost:9030/login' \
--header 'Content-Type: application/json' \
--data '{
    "email": "go@com.com",
    "password": "Password1~"
}'
```

- 取得商品推薦
  在 header 中加入 `Authorization` 並帶入登入後取得的 token，並且帶上 `Bearer` 字串

```shell
curl --location 'localhost:9030/products/recommendation' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwczovL2FsYW5jaGVuLmNvbSIsImV4cCI6MTcxMTE3ODk0NiwianRpIjoiNTNkOTBlYTktYmI2Ni00YjkwLWJkZjEtMjZkZjgyMmY3M2I3IiwiaWF0IjoxNzExMDkyNTQ2LCJpc3MiOiJBbGFuIGNoZW4iLCJuYmYiOjE3MTEwOTI1NDYsInN1YiI6IjI3MDA2YWE2LWI5NTctNDY0My05MTI5LTQ2NWNiNWYyNDZjYyJ9.QWIQQ3VcZkgJFLrskGDEJk4tAjKSE8RaxmkszBWnEdE'
```

## 效能測試

使用工具為 locust 測試，透過介面化的方式分析時間軸上的效能，觀察在快取失效的情況下，服務的效能表現

### 測試參數說明

為避免使用者因為 slow query 而等待，無法再下一秒持續發送請求，導致無法真實測試效能，，所以將每個使用者的任務數設定為 1，並且每秒新增使用者數量為 5，總共使用者數量為 600

### 測試結果

以下測試均在相同硬體環境下進行，為加速測試，將快取失效時間設定為 30 秒、測試時間為 2 分鐘

1. 使用快取來減少資料庫存取
   | 總請求數 | 失敗請求數 | 平均回應時間 | 最大回應時間 | 最小回應時間 | 90% 回應時間 | 95% 回應時間 | 99% 回應時間 |
   | -------- | ---------- | ------------ | ------------ | ------------ | ------------ | ------------ | ------------ |
   | 600 | 0 | 359.75ms | 3075ms | 3ms | 3000ms | 3000ms | 3000ms |
2. 透過 `single flight` 來減少重複的資料庫存取
   | 總請求數 | 失敗請求數 | 平均回應時間 | 最大回應時間 | 最小回應時間 | 90% 回應時間 | 95% 回應時間 | 99% 回應時間 |
   | -------- | ---------- | ------------ | ------------ | ------------ | ------------ | ------------ | ------------ |
   | 600 | 0 | 209.92ms | 3020ms | 2ms | 1000ms | 2000ms | 3000ms |

### 本地測試

- 需要環境: Python 3.7
- 安裝說明
  ```shell
  pip install locust
  ```
- 使用方式
  - 先透過 curl 登入取得 token，並將 token 貼到 locustfile.py 中
  ```shell
  locust -f test/locustfile.py
  ```
  - 在瀏覽器中開啟 `http://localhost:8089` 並設定使用者數量、Ramp up 等參數
- 參數設定
  - 總共使用者: 600
  - 每秒新增使用者(Ramp up): 5
  - 每個使用者執行的任務數: 1 (已經寫在 locustfile.py 中)

### TODO

- [ ] 當 redis 資料失效時，避免併發的請求同時向資料庫取得資料
      密碼從 gcp secret manager 等服務中取得
