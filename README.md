# TranslateBot
A Simple Telegram Translate Bot

第一个Go项目，写的代码贼垃圾

顺便练习JBIDE

准备写一篇GO学习日记发在博客上

## Usage
在Release中下载最新版本的可执行文件

上传至VPS，在同目录下创建 `bot.cfg`文件，写入以下内容
protocol支持 `http` `https`  `socks5`等协议
```json
{
    "token": "Your_Token",
    "proxy": {
		"enable": "no",
		"protocol": "socks5",
		"ip": "127.0.0.1",
		"port": 1080
	}
}
```
启动可执行文件

Linux: 
```shell script
screen -S tlbot
./tlbot
Ctrl A+D
```
Windows: 

双击`tlbot.exe`文件