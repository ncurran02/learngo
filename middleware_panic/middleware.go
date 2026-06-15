package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type RecoverMiddleware struct {
	Next http.Handler
}

func (m RecoverMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		dev, err := strconv.ParseBool(os.Getenv("DEV"))
		if err != nil {
			log.Fatalln("Unable to parse dev env")
		}

		if err1 := recover(); err1 != nil {
			log.Printf("\npanic recovered: %s", err1)

			var message string
			if dev {
				message = fmt.Sprintf("Internal Server Error\n%s", err1)
			} else {
				message = "Internal Server Error"
			}

			http.Error(w, message, http.StatusInternalServerError)
		}
	}()

	m.Next.ServeHTTP(w, r)
}
