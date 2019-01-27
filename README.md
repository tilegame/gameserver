# Tile Game Server

Repository for the tile game server.

[![Build Status](https://travis-ci.org/tilegame/gameserver.svg?branch=master)](https://travis-ci.org/tilegame/gameserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/tilegame/gameserver)](https://goreportcard.com/report/github.com/tilegame/gameserver)
[![GoDoc](https://godoc.org/github.com/tilegame/gameserver?status.svg)](https://godoc.org/github.com/tilegame/gameserver)



## Endpoints

Endpoint	| Description
--|--
https://thebachend.com/     | shows realtime information about endpoints
https://thebachend.com/cookie	    | generates and/or validates new cookies for clients
https://thebachend.com/scoreboard	| simple scoreboard example
https://thebachend.com/sessions	| generates a list of active sessions
wss://thebachend.com/ws	Main    | websocket connection for game (not implemented yet)
wss://thebachend.com/ws/echo	| echo server used for testing connection speeds


## Made With

- [The Go Programming Language](https://golang.org/).
- [Google Cloud's Compute Engine](https://cloud.google.com/compute/) for hosting the server.
- [Let's Encrypt](https://letsencrypt.org/) for the automatic free SSL/TLS Certificates
- [Gorilla web toolkit](http://www.gorillatoolkit.org/) for it's [gorilla/websocket](https://github.com/gorilla/websocket/) package.
