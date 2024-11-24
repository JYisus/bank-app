# Bank API

## Run the code

To execute the code you have the following options:

1. Go run (within project's root directory):

```bash
go run main.go
```

2. Compile the program and execute it (within project's root directory):

```bash
go build -o app main.go
./app
```

3. Docker (within project's root directory):

```bash
# Build the image
docker build -t bank-app .

# Run the image
docker run -o 8080:8080 bank-app
```

## Request examples

### Create account (POST /accounts)

This request will return the account id generated.

```bash
curl -X POST http://localhost:8080/accounts -H 'Content-Type: application/json' --data-raw '{"owner": "test", "initial_balance": 20}'
# {"id":"4ea66c9b-71f0-4502-9eb2-f77e0c3f3df2"}
```

### Get account (GET /accounts/{id})

```bash
curl -X GET http://localhost:8080/accounts/4ea66c9b-71f0-4502-9eb2-f77e0c3f3df2
# {"id":"0e9c52d5-138b-4437-a00a-c78503ffbc70","owner":"test","balance":20}
```

### List all account (GET /accounts)

```bash
curl -X GET http://localhost:8080/accounts
# [{"id":"fcfcc0b5-64bb-4a6c-b802-3460cf8b3622","owner":"test","balance":20},{"id":"4ea66c9b-71f0-4502-9eb2-f77e0c3f3df2","owner":"test","balance":20}]
```

### Create transaction (POST /accounts/{id}/transactions)

This request will return the transaction id generated.

```bash
curl -X POST "http://localhost:8080/accounts/fcfcc0b5-64bb-4a6c-b802-3460cf8b3622/transactions" -H 'Content-Type: application/json' --data-raw '{"type": "deposit","amount": 20.3}'
# {"id":"fe8442b3-6a0c-4074-af3d-de51e8f47f68"}
```

### Retrieve transactions for an Account (GET /accounts/{id}/transactions)

```bash
curl -X GET "http://localhost:8080/accounts/fe8442b3-6a0c-4074-af3d-de51e8f47f68/transactions"
# [{"id":"fe8442b3-6a0c-4074-af3d-de51e8f47f68","accountId":"fcfcc0b5-64bb-4a6c-b802-3460cf8b3622","type":"deposit","amount":20.3,"timestamp":"2024-11-24T03:26:51.835490418Z"}]
```

### Retrieve transactions for an Account (GET /accounts/{id}/transactions)

```bash
curl -X GET "http://localhost:8080/accounts/fe8442b3-6a0c-4074-af3d-de51e8f47f68/transactions"
# [{"id":"fe8442b3-6a0c-4074-af3d-de51e8f47f68","accountId":"fcfcc0b5-64bb-4a6c-b802-3460cf8b3622","type":"deposit","amount":20.3,"timestamp":"2024-11-24T03:26:51.835490418Z"}]
```

### Transfer (POST /transfer)

```bash
curl -X POST "http://localhost:8080/transfer" -H 'Content-Type: application/json' --data-raw '{"from_account_id": "fcfcc0b5-64bb-4a6c-b802-3460cf8b3622","to_account_id": "4ea66c9b-71f0-4502-9eb2-f77e0c3f3df2","amount": 10}'
# {"from_account_id":"fcfcc0b5-64bb-4a6c-b802-3460cf8b3622","to_account_id":"4ea66c9b-71f0-4502-9eb2-f77e0c3f3df2","amount":10}
```

## Run the tests

To run the tests, execute the following command within project's root directory:

```bash
go test ./...
```

## Future improvements

- Including a DB repository
- Add configurable log level and server port
- Rollback mechanism to avoid data inconsistency with repositories
- Avoid using float32 for money
- Improve logging
