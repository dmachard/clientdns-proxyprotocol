# DNS Client with Proxy Protocol v2  

This is a Go-based DNS client that supports **Proxy Protocol v2**, allowing you to send DNS queries while preserving client source IP information. It utilizes the [`miekg/dns`](https://github.com/miekg/dns) and [`pires/go-proxyproto`](https://github.com/pires/go-proxyproto) libraries.  

## Features  
✅ Sends DNS queries using UDP  
✅ Supports **Proxy Protocol v2** to include source IP details  
✅ Allows setting **Type-Length-Value (TLV)** fields in the Proxy Protocol header  
✅ Supports DNS query types: **A, AAAA, MX**  

## Installation  

Make sure you have **Go** installed (>=1.18).  

1. Clone the repository:  
    ```sh
    git clone https://github.com/yourusername/dns-proxy-client.git
    cd dns-proxy-client
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Build the binary:

    ```sh
    go build -o dns-client
    ```

## Usage

Run the client with the following syntax:

```sh
./dns-client <dns_server> <port> <domain> <type(A,AAAA,MX)> [kv_key=kv_value]
```

Example to query an AAAA record with Proxy Protocol and a TLV field:

```sh
./dns-client 1.1.1.1 53 example.com AAAA 1=custom-data
```

## Dependencies
- miekg/dns – DNS library for Go
- pires/go-proxyproto – Proxy Protocol implementation