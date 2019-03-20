# Falcon-WS

Light-weight socks5 compatible proxy service kit using WebSocket.

*Notice: This project is just for study and fun.*

## How it works?
![preview](https://raw.githubusercontent.com/SUCHMOKUO/falcon-ws/master/doc/images/falcon-ws.png)

Falcon-WS use *Websocket* as it's transport layer protocal, So it can break some limitation set by the firewall (eg. only http/https). With the help of Websocket, you will be able to establish a full-duplex connection with the server through the firewall, and then you can access any service based on TCP on the internet. 

## Build

```
go get -u -v github.com/SUCHMOKUO/falcon-ws
# Then build falcon-ws-server and falcon-ws-client.
```

## Usage

```
falcon-ws-server --help
falcon-ws-client --help
```