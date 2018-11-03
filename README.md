# Tile Game Server

Repository for the tile game server.

[![Build Status](https://travis-ci.org/tilegame/gameserver.svg?branch=master)](https://travis-ci.org/tilegame/gameserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/tilegame/gameserver)](https://goreportcard.com/report/github.com/tilegame/gameserver)
[![GoDoc](https://godoc.org/github.com/tilegame/gameserver?status.svg)](https://godoc.org/github.com/tilegame/gameserver)



## Endpoints

URL | Dev Status | Description
----|------------|---------------
https://thebachend.com/ | Live | Shows realtime info about endpoints.
https://thebachend.com/login | planned | Login screen; HTTP POST to get a token.
wss://thebachend.com/ws | planned | Main websocket endpoint for the game
wss://thebachend.com/ws/echo | Live | Where the game currently is.



## Made With

- [The Go Programming Language](https://golang.org/).
- [Google Cloud's Compute Engine](https://cloud.google.com/compute/) for hosting the server.
- [Let's Encrypt](https://letsencrypt.org/) for the automatic free SSL/TLS Certificates
- [Gorilla web toolkit](http://www.gorillatoolkit.org/) for it's [gorilla/websocket](https://github.com/gorilla/websocket/) package.
