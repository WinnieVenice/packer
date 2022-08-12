package cq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/WinnieVenice/packer/model"
)

var (
	client = &http.Client{
		Timeout: time.Minute,
	}
)

func SendGroupMsg(groupId int64, msg string, autoEscape bool) error {
	body := map[string]interface{}{
		"group_id":    groupId,
		"message":     msg,
		"auto_escape": autoEscape,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("packer marshal http body failed, err = (%s)\n", err.Error())
		return err
	}

	url := fmt.Sprintf("%s/%s", model.UrlRobotBase, model.UrlSendGroupMsg)
	httpReq, err := http.NewRequest(model.MethodSendGroupMsg, url, bytes.NewReader(jsonBody))
	httpReq.Header.Add("Content-Type", model.ContentTypeSendGroupMsg)
	httpReq.Header.Add("X-SELF-ID", model.SelfId)

	if err != nil {
		fmt.Printf("packer new http req failed, err = (%s)\n", err.Error())
		return err
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("packer send http failed, err = (%s)\n", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("packer send http failed, retCode = (%d)\n", resp.StatusCode)
		return err
	}
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("packer get http resp body failed, err = (%s)\n", err.Error())
		return err
	}
	fmt.Println(string(jsonData))

	return nil
}

func MSendGroupMsg(groupIds []int64, msg string, autoEscape bool) {
	wg := sync.WaitGroup{}
	for _, id := range groupIds {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			err := SendGroupMsg(i, msg, autoEscape)
			if err != nil {
				return
			}
		}(id)
	}
	wg.Wait()
}
