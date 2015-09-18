# sirbot

sirbot is "Slack IRC relay bot"

```
Usage of sirbot:
  -irc-channel string
    	IRC channel to join
  -irc-host string
    	IRC server host (default "localhost")
  -irc-password string
    	IRC server password
  -irc-port int
    	IRC server port (default 6666)
  -irc-secure
    	IRC use SSL
  -listen string
    	HTTP listen address (for accept Outgoing Webhook) (default ":7777")
  -nick string
    	IRC nick (default "sirbot")
  -slack-channel string
    	Slack channel to join
  -webhook-token string
    	Slack Outgoing Webhook token
  -webhook-url string
    	Slack Incomming Webhook URL
```

## Install

```
go get github.com/fujiwara/sirbot
```

## Usage example

```
$ sirbot \
  -irc-host irc.example.com \
  -irc-port 6666 \
  -irc-channel "#sandbox" \
  -irc-password secret \
  -listen :7777 \
  -slack-channel "#sandbox" \
  -webhook-token xxxxxxxxxxxx \
  -webhook-url https://hooks.slack.com/services/xxxx/yyyy/zzzzzzzzz \
```

## Author

2015 Fujiwara Shunichiro.

## License

MIT License

