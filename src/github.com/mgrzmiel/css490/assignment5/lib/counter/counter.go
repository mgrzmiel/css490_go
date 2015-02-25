package counter

// import(
// 	"github.com/mgrzmiel/css490/assignment5/lib/counter"
// )

const (
	kindIncerement = iota //value 0
	kindGet
)

type request struct {
	// response changgel for Get requests
	resp chan int
	// flag to determine which operation is requested: increment
	// or get
	requestKind int
	// variable name
	key string
	// amount to increment (or decrement if negative)
	delta int
}

type Counter struct {
	req chan request
}

// New creates a Counter object
// A counter is a map between name and value, where each value can be
// incremented atomically.
//
// There is also an unimplemented Dump function to return a copy of
// the dictionary (for exporting to another service..)
func New() *Counter {
	// A counter object just holds a (private) channel for
	// communication with the goroutine service.
	c := &Counter{
		req: make(chan request),
	}
	// spawn off the service goroutine, passing the counter
	// object to establish communication
	go counter(c)
	// And return the counter object.
	return c
}

// Get returns the value of the given counter object.  It provides a
// conventional interface, so the calling function does not need to
// know that this is implemented by communicating with another
// thread.
func (c *Counter) Get(key string) int {
	// Create a channel to receive the response.
	resp := make(chan int)
	// Send a request to the service routine via the counter's
	// request field.
	c.req <- request{
		// Create channel to receive response
		resp: resp,
		// Operation is get, not increment
		requestKind: kindGet,

		// Variable to get value of.
		key: key,
	}
	// And return the response when you get it.
	return <-resp
}

// Incr increments (or decrements if delta is negative) a couunter.
// Again, we have a convtional interface that hides the fact that it
// is implemented by communication with another thread (goroutine).
func (c *Counter) Incr(key string, delta int) {
	c.req <- request{
		resp:        nil,
		requestKind: kindIncerement,
		key:         key,
		delta:       delta,
	}
}

// counter runs in its own goroutine to service Counter object
// requests.
func counter(c *Counter) {
	// private variable to hold the map
	data := make(map[string]int)
	// range on the channel means we can exit the loop and hence
	// terminate the thread by closing the channel
	for req := range c.req {
		switch req.requestKind {
		// Service the request
		//here i have sth wrong
		case kindIncerement: //req.increment:
			data[req.key] = data[req.key] + req.delta
		case kindGet:
			req.resp <- data[req.key]
		}

	}
}
