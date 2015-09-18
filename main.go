package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	irc "github.com/thoj/go-ircevent"
)

const Slackbot = "slackbot"

var (
	WebhookURL   string
	WebhookToken string
	ListenAddr   string
	SlackChannel string
	IRCChannel   string
	IRCNick      string
	IRCHost      string
	IRCPort      int
	IRCPassword  string
	IRCSecure    bool
)

func main() {
	flag.StringVar(&WebhookURL, "webhook-url", "", "Slack Incomming Webhook URL")
	flag.StringVar(&WebhookToken, "webhook-token", "", "Slack Outgoing Webhook token")
	flag.StringVar(&ListenAddr, "listen", ":7777", "HTTP listen address (for accept Outgoing Webhook)")
	flag.StringVar(&IRCNick, "nick", "sirbot", "IRC nick")
	flag.StringVar(&IRCHost, "irc-host", "localhost", "IRC server host")
	flag.IntVar(&IRCPort, "irc-port", 6666, "IRC server port")
	flag.StringVar(&IRCPassword, "irc-password", "", "IRC server password")
	flag.BoolVar(&IRCSecure, "irc-secure", false, "IRC use SSL")
	flag.StringVar(&IRCChannel, "irc-channel", "", "IRC channel to join")
	flag.StringVar(&SlackChannel, "slack-channel", "", "Slack channel to join")
	flag.Parse()

	ch := make(chan Message, 10)
	go startHttpServer(ch)

	slack := &SlackAgent{
		WebhookURL: WebhookURL,
		client:     &http.Client{},
	}
	agent := irc.IRC(IRCNick, IRCNick)
	agent.UseTLS = IRCSecure
	agent.Password = IRCPassword
	err := agent.Connect(fmt.Sprintf("%s:%d", IRCHost, IRCPort))
	if err != nil {
		log.Println(err)
		return
	}

	agent.AddCallback("001", func(e *irc.Event) {
		agent.Join(IRCChannel)
		log.Println("Joined", IRCChannel)
	})

	agent.AddCallback("PRIVMSG", func(e *irc.Event) {
		msg := Message{
			Channel: SlackChannel,
			Text:    e.Message(),
			UserName: fmt.Sprintf(
				"%s[%s]",
				e.Nick,
				e.Arguments[0], // channel
			),
		}
		err := slack.Post(msg)
		if err != nil {
			log.Println("[error]", err)
		}
	})
	go func() {
		for {
			msg := <-ch
			agent.Privmsg(IRCChannel, msg.Text)
			time.Sleep(1 * time.Second)
		}
	}()

	agent.Loop()
}

func startHttpServer(ch chan Message) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if token := req.FormValue("token"); token != WebhookToken {
			log.Println("[warn] Invalid token:", token)
			return
		}
		userName := req.FormValue("user_name")
		if userName == Slackbot {
			return
		}
		text := fmt.Sprintf("[%s] %s", userName, req.FormValue("text"))
		msg := Message{
			Channel:  IRCChannel,
			UserName: userName,
			Text:     text,
		}
		select {
		case ch <- msg:
		default:
		}
	})
	log.Fatal(http.ListenAndServe(ListenAddr, nil))
}

type Message struct {
	Channel   string `json:"channel"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
	UserName  string `json:"username"`
}

type SlackAgent struct {
	WebhookURL string
	client     *http.Client
}

func (a *SlackAgent) Post(m Message) error {
	payload, _ := json.Marshal(&m)
	v := url.Values{}
	v.Set("payload", string(payload))
	log.Println("post to slack", a, string(payload))
	resp, err := a.client.PostForm(a.WebhookURL, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		return fmt.Errorf("failed post to slack:%s", body)
	} else {
		return err
	}
}
