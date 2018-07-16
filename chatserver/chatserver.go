package chatserver

type client chan<-string

var (
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string)
	clientMap = make(map[client]bool)
)

func broadcaster() {
	for {
		select {
		case m <- messages:
			for c := range clientMap {
				c <- m
			}
		case c <- entering:
			clientMap[c] = true
		case c <- leaving:
			delete(clientMap, c)
			close(c)
		}
	}
}

func handleChat() {

}
