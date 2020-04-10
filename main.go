package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
)

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		DB:              0,
		MaxRetries:      5,
		MaxRetryBackoff: time.Second * 1,
	})

	conn := rc.Conn()
	conn.Select(0)
	http.HandleFunc("/getdata", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "Method Not Allowed")
			return
		}

		conn := rc.Get("mykey")
		if conn.Err() != nil {
			fmt.Println(conn.Err().Error())
			io.WriteString(w, conn.Err().Error())
			return
		}

		data, _ := conn.Result()

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, data)
		return
	})

	http.ListenAndServe(":8080", nil)
}
