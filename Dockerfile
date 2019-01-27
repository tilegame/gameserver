FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /go/src/github.com/tilegame/gameserver

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# ENV PORT 8080

CMD ["gameserver", "-a", ":80", "tls"]

EXPOSE 443
