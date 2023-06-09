package util

import (
	"fmt"
	"net/http"
)

var responde struct {
	status int
	msg    string
}

func Contains(indices []int, indice int) bool {
	for _, s := range indices {
		if s == indice {
			return true
		}
	}
	return false

}

func JSON(w http.ResponseWriter, statusCode int, code, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintln(code)))
	w.Write([]byte(fmt.Sprintln(message)))

}
