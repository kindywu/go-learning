package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

// SSEHandler handles SSE requests
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Set the headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a flusher to manually flush data to the client
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	name := rand.Int()
	fmt.Printf("sse client connected %d \n", name)
	// Create a channel to signal when the client disconnects
	notify := r.Context().Done()

	go func() {
		<-notify
		fmt.Printf("sse client disconnected %d \n", name)
	}()

	// Send some data to the client every second
	for {
		select {
		case <-notify:
			return
		default:
			fmt.Fprintf(w, "data: %s\n\n", time.Now().String())
			flusher.Flush()
			time.Sleep(1 * time.Second)
		}
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SSE Example</title>
</head>
<body>
    <h1>SSE Example</h1>
    <p>Open Console</p>
    <script>
        const es = new EventSource("http://localhost:8080/events");
        es.onopen = () => console.log("Connection Open!");
        es.onmessage = (e) => console.log("Message:", e.data);
        es.onerror = (e) => {
            console.log("Error:", e);
            // es.close();
        };
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/events", SSEHandler)

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
