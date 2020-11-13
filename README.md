# keva_ipfs

This is the backend for the Keva mobile app to upload and pin IPFS files.

Environment variables to set: `KEVA_PAYMENT_ADDRESS` and `KEVA_MIN_PAYMENT`. E.g.

```
export KEVA_PAYMENT_ADDRESS=VJnZF9xkqQPcgceyq2FXvmsg5Lmb9UUGyi
export KEVA_MIN_PAYMENT=100
```

This server will only pin an IPFS file if a payment is made to `KEVA_PAYMENT_ADDRESS`, with the minimal Kevacoin amount defined by `KEVA_MIN_PAYMENT`.




