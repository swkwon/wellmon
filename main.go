package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"wellmon/date"

	"github.com/antchfx/htmlquery"
)

func getHTML(URL string) {
	doc, err := htmlquery.LoadURL(URL)
	if err != nil {
		log.Fatal(err)
	}

	nodes := htmlquery.Find(doc, "//a/@class")
	if len(nodes) <= 0 {
		log.Println("node count = 0", URL)
	}
	for _, node := range nodes {
		if "btn_buy fe" == htmlquery.InnerText(node) {
			titles := htmlquery.Find(doc, "//title")
			for _, t := range titles {
				title := htmlquery.InnerText(t)
				noti(URL, title)
				return
			}
		}
	}
}

func send(text string) error {
	telegramURLFormat := "https://api.telegram.org/bot%s/SendMessage"
	channel := os.Getenv("WK_CHANNEL")
	token := os.Getenv("WK_TOKEN")
	url := fmt.Sprintf(telegramURLFormat, token)
	msg := fmt.Sprintf(`{"chat_id":"%s","text":"%s"}`, channel, text)

	reader := bytes.NewBuffer([]byte(msg))
	var err error
	if req, err := http.NewRequest("POST", url, reader); err == nil {
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		if res, err := client.Do(req); err == nil {
			defer res.Body.Close()
			_, e := ioutil.ReadAll(res.Body)
			if e != nil {
				return e
			}
		} else {
			return err
		}
	}
	return err
}

func noti(URL, title string) {
	err := send("떴다! " + title + " " + URL)
	if err != nil {
		log.Println(err)
	}
	date.Save(URL)
}

func check(URL string) bool {
	v := date.Load(URL)
	if v == "" {
		return false
	}

	t, e := time.Parse(time.RFC3339, v)
	if e != nil {
		return false
	}

	now := date.GetSeoulTime()
	if now.Year() == t.Year() &&
		now.Month() == t.Month() &&
		now.Day() == t.Day() &&
		now.Hour() == t.Hour() {
		return true
	}

	return false
}

func main() {

	URLs := []string{
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=1007206&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=bWV9",
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=1007205&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=a253UA%3D%3D",
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=922816&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=bmt8W14%3D",
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=920693&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=aWp3Ug%3D%3D",
	}

	for _, v := range URLs {
		if skip := check(v); skip == false {
			getHTML(v)
			time.Sleep(1 * time.Second)
		}
	}
}
