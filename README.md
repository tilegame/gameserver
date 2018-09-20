# Ninja Server

[![Build Status](https://travis-ci.org/fractalbach/ninjaServer.svg?branch=master)](https://travis-ci.org/fractalbach/ninjaServer)
[![Go Report Card](https://goreportcard.com/badge/github.com/fractalbach/ninjaServer)](https://goreportcard.com/report/github.com/fractalbach/ninjaServer)
[![GoDoc](https://godoc.org/github.com/fractalbach/ninjaServer?status.svg)](https://godoc.org/github.com/fractalbach/ninjaServer)

Repository for the server of the [Tile Experiments Project](https://github.com/fractalbach/TileExperiments)
(Still in early development).


## Endpoints

URL | Dev Status | Description
----|------------|---------------
https://thebachend.com/ | Live | Shows realtime info about endpoints.
https://thebachend.com/login | planned | Login screen; HTTP POST to get a token.
wss://thebachend.com/ws | planned | Main websocket endpoint for the game
wss://thebachend.com/ws/echo | Live | Where the game currently is.



## Credits 

- [Google Cloud's Compute Engine](https://cloud.google.com/compute/) for hosting the server.
- [Let's Encrypt](https://letsencrypt.org/) for the automatic free SSL/TLS Certificates
- [Gorilla web toolkit](http://www.gorillatoolkit.org/) for it's [gorilla/websocket](https://github.com/gorilla/websocket/) package.
