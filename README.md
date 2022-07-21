# remote-signing-api

## Nonce Manager

Minimal implementation of nonce manager to make sure transactions are properly numbered. It is a task of the wallet/signer to assigne correct nonces to transactions since it is not possible to rely on the information from the node.

The Nonce Manager allows to track nonces for a given address and a chain ID


### How to Run

```
go run server/server.go config.json
```

### API

```
GET /nonce/:chainId/:address
```

Returns a `NonceResponse` object in response. It contains the nonce and the account information. This call initializes the nonce value from blockchain if new address is provided and internally tracks and increases the nonce value on each call.

It is a responsibility of client to *return* the nonce in case of the TX is not added to the blockchain

```
PUT /nonce/:chainId/:address/:nonce
```

Puts an unused `nonce` back to the Nonce Manager - this needs to be called in case of errors on transaction submission, timeouts etc.

```

```
POST /nonce/:chainId/:address/sync
```

Synchronizes the nonce value with the blockchain. Cleans returned nonces list.