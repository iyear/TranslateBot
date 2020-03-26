package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	var (
		BotToken string
		Socks5   string
	)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Read Config ERROR")
	}
	BotToken = viper.GetString("bot_token")
	fmt.Println("Token:" + BotToken)
	Socks5 = viper.GetString("socks5")
	botsettings := tb.Settings{
		Token:  BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}
	if Socks5 != "" {
		fmt.Println("Proxy:" + Socks5)
		dialer, err := proxy.SOCKS5("tcp", Socks5, nil, proxy.Direct)
		if err != nil {
			log.Fatal("Error creating dialer, aborting.")
		}
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
		botsettings.Client = httpClient
	}

	b, err := tb.NewBot(botsettings)
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "欢迎使用Translatebot!")
	})
	var btn1 = tb.InlineButton{
		Unique: "ZH_CN2EN",
		Text:   "中文->英文",
	}
	var btn2 = tb.InlineButton{
		Unique: "ZH_CN2JA",
		Text:   "中文->日文",
	}
	var inlinekeys = [][]tb.InlineButton{{btn1}, {btn2}}
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
	resptxt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR!")
		return ""
	}
	value := gjson.Get(string(resptxt), "translateResult.0.0.tgt")
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
