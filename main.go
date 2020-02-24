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
	wmlog "wellmon/log"

	"github.com/antchfx/htmlquery"
)

func getHTML(URL string) {
	found := false
	defer func() {
		if found == false {
			wmlog.DLog("button not found", URL)
		}
	}()
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
			log.Println("!!! OPEN !!!", URL)
			titles := htmlquery.Find(doc, "//title")
			for _, t := range titles {
				title := htmlquery.InnerText(t)
				noti(URL, title)
				found = true
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
	now := date.GetSeoulTime()
	if now.Hour() < 9 || now.Hour() > 17 {
		return true
	}

	v := date.Load(URL)
	if v == "" {
		return false
	}

	t, e := time.Parse(time.RFC3339, v)
	if e != nil {
		return false
	}

	if now.Year() == t.Year() &&
		now.Month() == t.Month() &&
		now.Day() == t.Day() &&
		now.Hour() == t.Hour() {
		return true
	}

	return false
}

func processingPassword(v string) string {
	token := os.Getenv("WK_TOKEN")
	if len(token) <= 0 {
		return "no token"
	} else {
		var ret string
		for i := 0; i < len(token); i++ {
			ret = ret + "*"
		}
		return ret
	}
}

func checkEnv() {
	log.Println("check channel... : ", os.Getenv("WK_CHANNEL"))
	log.Println("check token... : ", processingPassword(os.Getenv("WK_TOKEN")))
}

var hour int

func logPerHour() {
	now := date.GetSeoulTime()
	if hour != now.Hour() {
		hour = now.Hour()
		log.Println(now)
	}
}

func main() {
	hour = -1
	log.Println("start wellmon...")
	checkEnv()
	URLs := []string{
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=1007206&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=bWV9",
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=1007205&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=a253UA%3D%3D",
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=922816&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=bmt8W14%3D",
		"http://www.welkeepsmall.com/shop/shopdetail.html?branduid=920693&xcode=023&mcode=001&scode=&type=X&sort=regdate&cur_code=023001&GfDT=aWp3Ug%3D%3D",
	}

	for {
		logPerHour()
		for _, v := range URLs {
			if skip := check(v); skip == false {
				getHTML(v)
				time.Sleep(1 * time.Second)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
