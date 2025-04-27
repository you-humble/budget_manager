# BUDGET MANAGER (IN PROCESS)

An application for managing your incomes and expenses.
It can also tracking bank deposits and a rate of the most popular criptocurrencies 

## Run
go run main.go

---

## Running the tests
go test -v

---
## Expected data
Create wallet

POST http://localhost:8080/wallet/create

```
body:
{
  "user_id": 1,
  "title": "my wallet",
  "general": 25000
}
```

---

Show wallet

GET http://localhost:8080/wallet/show
```
body:
{
    "user_id": 1
}
```
---

Add operation

POST http://localhost:8080/wallet/operation/add
```
body:
{
  "user_id": 1,
  "operation": {
    "wallet_id": 1,
    "type": "income",
    "amount": 10000
  }
}
```
---

{
  "login": "user",
  "password": "123"
}