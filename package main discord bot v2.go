package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + "Discord token")
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	defer discord.Close()

	fmt.Println("NASA Apod is now running. CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case "hi":
		s.ChannelMessageSend(m.ChannelID, "heyo")
	case "!nasa apod":
		res, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=NasaApiKey")
		if err != nil {
			fmt.Println("Error fetching APOD:", err)
			return
		}
		defer res.Body.Close()

		var result map[string]interface{}
		responseData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading response data:", err)
			return
		}

		if err := json.Unmarshal(responseData, &result); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, "**Date**: "+result["date"].(string))
		s.ChannelMessageSend(m.ChannelID, "**Title**: "+result["title"].(string))
		s.ChannelMessageSend(m.ChannelID, "**Explanation**: \n"+result["explanation"].(string))
		s.ChannelMessageSend(m.ChannelID, result["url"].(string))
	default:
		fmt.Println(m.Content)
	}
}