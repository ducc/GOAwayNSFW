package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"github.com/koyachi/go-nude"
	"encoding/json"
)

type config struct {
	Token string `json:"token"`
}

var (
	botId string
)

func loadConfig() config {
	var conf config
	body, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(body, &conf)
	return conf
}

func main() {
	conf := loadConfig()
	discord, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Fatalln(err)
	}
	user, err := discord.User("@me")
	if err != nil {
		log.Fatalln(err)
	}
	botId = user.ID
	discord.AddHandler(messageCreate)
	err = discord.Open()
	if err != nil {
		log.Fatalln(err)
	}
	<-make(chan struct{})
}

func saveImage(url string) {
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fileName := resp.Request.URL.Path
	err = ioutil.WriteFile("img/" + fileName[strings.LastIndex(fileName, "/")+1:], body, os.ModeAppend)
	if err != nil {
		log.Fatalln(err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botId {
		return
	}
	log.Println("got msg:", m.Content)
	if len(m.Attachments) == 0 {
		return
	}
	for _, attachment := range m.Attachments {
		saveImage(attachment.URL)
		res, err := nude.IsNude("img/" + attachment.Filename)
		if err != nil {
			log.Fatalln(err)
		}
		if res {
			s.ChannelMessageSend(m.ChannelID, "porn!!! ban!!!")
		}
		os.Remove("img/" + attachment.Filename)
	}
}