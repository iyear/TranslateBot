# TranslateBot
A Simple Telegram Translate Bot

第一个Go项目，写的代码贼垃圾

顺便练习JBIDE

准备写一篇GO学习日记发在博客上

## Usage
在Release中下载最新版本的可执行文件

学会了yaml+viper所以改成这种方式

上传至VPS，在同目录下创建 `config.yml`文件，写入以下内容
```yaml
bot_token: xxxxx
socks5: 127.0.0.1:1080
#如果不需要代理删去此行
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