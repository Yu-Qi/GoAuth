# GoAuth

提供簡易的帳號註冊、驗證、登入等服務

### 所需軟體/服務

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
