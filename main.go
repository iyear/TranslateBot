package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	os.Setenv("HTTPS_PROXY", "socks5://127.0.0.1:1080")
	os.Setenv("HTTP_PROXY", "socks5://127.0.0.1:1080")
	botcfg, err := ioutil.ReadFile("./bot.cfg")
	if err != nil {
		fmt.Println("读取token失败:", err)
	}
	token := gjson.Get(string(botcfg), "token")
	if token.String() != "" {
		fmt.Println("token:" + token.String())
	}
	b, err := tb.NewBot(tb.Settings{
		Token:  token.String(),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "欢迎使用Translatebot!")
	})
	btn1 := tb.InlineButton{
		Unique: "ZH_CN2EN",
		Text:   "中文->英文",
	}
	btn2 := tb.InlineButton{
		Unique: "ZH_CN2JA",
		Text:   "中文->日文",
	}
	inlinekeys := [][]tb.InlineButton{
		[]tb.InlineButton{btn1},
		[]tb.InlineButton{btn2},
	}
	b.Handle(&btn1, func(c *tb.Callback) {
		b.Respond(c)
		fmt.Println(btn1.Unique)
		b.Reply(c.Message, GetBetweenStr(c.Message.Text, "$", "$")+"  -->  "+translate(GetBetweenStr(c.Message.Text, "$", "$"), btn1.Unique))
	})
	b.Handle(&btn2, func(c *tb.Callback) {
		b.Respond(c)
		fmt.Println(btn2.Unique)
		b.Reply(c.Message, GetBetweenStr(c.Message.Text, "$", "$")+"  -->  "+translate(GetBetweenStr(c.Message.Text, "$", "$"), btn2.Unique))
	})
	b.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Println(m.Text)
		b.Send(m.Sender, "开始翻译，请选择! $"+m.Text+"$", &tb.ReplyMarkup{
			InlineKeyboard: inlinekeys,
		})
	})

	b.Start()
}
func translate(s, trtype string) string {
	resp, err := http.Get("http://fanyi.youdao.com/translate?&doctype=json&type=" + trtype + "&i=" + s)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	zh_CN, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR!")
		return ""
	}
	value := gjson.Get(string(zh_CN), "translateResult.0.0.tgt")
	return value.String()
}
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
