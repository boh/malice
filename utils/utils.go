package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/dustin/go-jsonpointer"
	"github.com/jordan-wright/email"
	er "github.com/maliceio/malice/malice/errors"
)

// ParseJSON returns a JSON value for a given key
// NOTE: https://godoc.org/github.com/dustin/go-jsonpointer
func ParseJSON(data []byte, path string) (out string) {
	var o map[string]interface{}
	er.CheckError(json.Unmarshal(data, &o))
	out = jsonpointer.Get(o, string(data)).(string)
	return
}

// ParseMail takes in an HTTP Request and returns an Email object
// TODO: This function will likely be changed to take in a []byte
func ParseMail(r *http.Request) (email.Email, error) {
	e := email.Email{}
	m, err := mail.ReadMessage(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(m.Body)
	e.HTML = body
	return e, err
}

// AskForConfirmation prompts user for yes/no response
func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "yes"}
	nokayResponses := []string{"n", "no"}
	if stringInSlice(strings.ToLower(response), okayResponses) {
		return true
	}
	if stringInSlice(strings.ToLower(response), nokayResponses) {
		return false
	}
	fmt.Println("Please type yes or no and then press enter:")
	return AskForConfirmation()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
