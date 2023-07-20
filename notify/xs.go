package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type typeText struct {
	At      notifyAt        `json:"at"`
	Text    typeTextContent `json:"text"`
	Msgtype string          `json:"msgtype"`
}

type typeTextContent struct {
	Content string `json:"content"`
}

type typeMarkdown struct {
	At       notifyAt            `json:"at"`
	Markdown typeMarkdownContent `json:"markdown"`
	Msgtype  string              `json:"msgtype"`
}

type typeMarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type notifyAt struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type notifyResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

const httpTimout = 10 * time.Second

type notifyXs struct {
	webhook string
}

func (n *notifyXs) Name() string {
	return "xs"
}

func (n *notifyXs) Webhook() string {
	return n.webhook
}

func (n *notifyXs) SetWebHook(ctx context.Context, webhook string) {
	n.webhook = webhook
}

// NotifyText 推送文本
func (n *notifyXs) NotifyText(ctx context.Context, content string, atUserIds []string) error {
	if n.webhook == "" {
		return nil
	}
	data := typeText{
		Msgtype: "text",
		Text:    typeTextContent{Content: content},
		At: notifyAt{
			AtUserIds: atUserIds,
		},
	}

	jsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := request(jsonByte, n.webhook); err != nil {
		return err
	}

	return nil
}

// NotifyMarkdown 推送markdown
func (n *notifyXs) NotifyMarkdown(ctx context.Context, content string, atUserIds []string) error {
	if n.webhook == "" {
		return nil
	}

	data := typeMarkdown{
		Msgtype: "markdown",
		Markdown: typeMarkdownContent{
			Title: "*熔断通知*",
			Text:  content,
		},
		At: notifyAt{
			AtUserIds: atUserIds,
		},
	}

	jsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := request(jsonByte, n.webhook); err != nil {
		return err
	}

	return nil
}

func request(message []byte, url string) (notifyResp, error) {
	var res notifyResp
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(message))
	if err != nil {
		return res, err
	}
	req.Header.Add("Accept-Charset", "utf8")
	req.Header.Add("Content-Type", "application/json")

	client := new(http.Client)
	client.Timeout = httpTimout
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resultByte, &res)
	if err != nil {
		return res, fmt.Errorf("unmarshal http response body from json error = %w", err)
	}

	if res.Errcode != 0 {
		return res, fmt.Errorf("send message error = %s", res.Errmsg)
	}

	return res, nil
}
