// Copyright © 2018 Charles Corbett <nafredy@gmail.com>
//

package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// FormatHostPort returns a string formatted host:port
func FormatHostPort(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// RespondJSON Formats a response as a JSON Payload
func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ToInt converts a string to an int
func ToInt(input string) int {
	i, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("%v", err)
	}

	return i
}
