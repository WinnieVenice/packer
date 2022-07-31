package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestMain(t *testing.T) {
	s := ""
	body, _ := json.Marshal(s)
	req, _ := http.NewRequest("POST", "http://localhost:5701/", bytes.NewReader(body))
	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println("send http err = ", err)
		return
	}
	fmt.Println("resp = ", resp)
}
