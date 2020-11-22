# keva_ipfs

This is the backend for the Keva mobile app to upload and pin IPFS files.

Environment variables to set:

```
# Connection to ElectrumX server
export KEVA_ELECTRUM_HOST=ec0.kevacoin.org
export KEVA_ELECTRUM_SSL_PORT=50002

# 5 KVA for uploading image
export KEVA_MIN_PAYMENT=5

# Payment address
export KEVA_PAYMENT_ADDRESS=VQfWnB3aUyzYfTt...

# TLS/SSL setting
export KEVA_TLS_ENABLED=1
export KEVA_TLS_KEY=/etc/letsencrypt/live/$KEVA_ELECTRUM_HOST/privkey.pem
export KEVA_TLS_CERT=/etc/letsencrypt/live/$KEVA_ELECTRUM_HOST/fullchain.pem
```

This server will only pin an IPFS file if a payment is made to `KEVA_PAYMENT_ADDRESS`, with the minimal Kevacoin amount defined by `KEVA_MIN_PAYMENT`. It assumes that TLS/SSL is used, and the certificate is issued by Letsencrypt.

The backend listens to port `$KEVA_ELECTRUM_SSL_PORT + 10`. E.g. if the ElectrumX server is listening on port 50002, this backend will listen to port 50012.




