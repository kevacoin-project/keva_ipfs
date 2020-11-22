# keva_ipfs

This is the backend for the Keva mobile app to upload and pin IPFS files.

## Prerequsites

It assumes that you already have an EletrumX server running on the same server.

### Download and install IPFS

You must run an IPFS peer on the server. To download and install an IPFS peer:

```
wget https://dist.ipfs.io/go-ipfs/v0.7.0/go-ipfs_v0.7.0_linux-amd64.tar.gz
tar -xvzf go-ipfs_v0.7.0_linux-amd64.tar.gz
cd go-ipfs
sudo bash install.sh
ipfs init   # Initialize the IPFS peer. Only need to run once.
ipfs daemon   # Start the daemon
```

### Download and install Golang

The backend is written in Golang and you need Golang to build the server:

```
wget https://golang.org/dl/go1.15.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

## Install the Server

Clone this repo and build the server:

```
git clone https://github.com/kevacoin-project/keva_ipfs
cd keva_ipfs
go build .
```

Environment variables to set:

```
# Connection to ElectrumX server
export KEVA_ELECTRUM_HOST=ec0.kevacoin.org   # Change it to your ElectrumX server!
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

After setting the environment variables, start the server:

```
./go_be
```
