# Metronero Backend
Backend and API server for Metronero.

Metronero is a checkout system for Monero cryptocurrency. It allows merchants to generate payment pages upon checkout. Instances have the ability to house multiple merchant accounts and be monetized with commission.

## Requirements
* [MoneroPay](https://gitlab.com/moneropay/moneropay)
* PostgreSQL server

## Installation
```sh
$ go install gitlab.com/metronero/backend 
```

## Usage
```
Usage of ./backend:
  -bind="localhost:5001": Bind address:port
  -callback-url="http://localhost:8080/callback": Incoming callback URL
  -log-format="pretty": Log format (pretty or json)
  -moneropay="http://localhost:5000": MoneroPay instance
  -postgres="postgresql://metronero:m3tr0n3r0@localhost:5432/metronero?sslmode=disable": PostgreSQL connection string
  -token-secret="aabbccddeeffgg": JWT secret

```

# API
The API client library in Go is located [here](https://gitlab.com/metronero/backend/-/tree/master/pkg/api).

## Admin
Endpoints for instance-wide queries or management. Requires admin privileges.

### GET /admin
Response data:
```json
{
    "instance_stats": {
        "wallet_balance": 10000000,
        "total_profits": 5000000000,
        "total_merchants": 12
    },
    "recent_withdrawals": [
        {
            "id": "009c1085-3988-4426-b655-fdf8add5a0d1",
            "merchant_name": "siren",
            "amount": 10000000,
            "date": "2023-04-08T19:19:05Z"
        },
        {
            "id": "658bfae2-8669-4639-8c6b-c2d8aa75a0c8",
            "merchant_name": "stnby",
            "amount": 10000000,
            "date": "2023-03-29T12:01:40Z"
        }
    ]
}
```

### GET /admin/withdrawal


### GET /admin/payment

### GET /admin/merchant

### GET /admin/merchant/{account_id}

### POST /admin/merchant/{account_id}

### DELETE /admin/merchant/{account_id}

### GET /admin/instance
Response body:
```json
{
    "version": "0.0.1",
    "default_commission": 1000,
    "custodial_mode": true,
    "registrations_allowed": true,
    "withdrawal_times": "weekly"
}
```

### POST /admin/instance

## Merchant
Endpoints for payment automation, merchant account/store management.

### GET /merchant

### GET /merchant/payment

### POST /merchant/payment

### GET /merchant/withdrawal

## Authentication
### POST /login
### POST /register

# Contributing
Issues and merge requests are only checked on [GitLab](https://gitlab.com/metronero/backend/).\
Alternatively, you can send patch files via email at [moneropay@kernal.eu](mailto:moneropay@kernal.eu).\
For development related discussions and questions join [#moneropay:kernal.eu](https://matrix.to/#/#moneropay:kernal.eu) Matrix group.

# Donations
Metronero does not receive funding from Monero's Community Crowdfunding System (CCS). We believe that CCS is a faulty/corrupt system that doesn't reflect the "community" in its name therefore, we don't recommend it. You can help the Metronero project by donating here:

- `46VGoe3bKWTNuJdwNjjr6oGHLVtV1c9QpXFP9M2P22bbZNU7aGmtuLe6PEDRAeoc3L7pSjfRHMmqpSF5M59eWemEQ2kwYuw`
- Or [https://donate.kernal.eu](https://donate.kernal.eu) - if you would like to leave a message! 