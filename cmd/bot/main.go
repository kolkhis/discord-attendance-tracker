package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Starting Discord bot...")
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable not set")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
		return
	}

	// Request intents
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates

	// Register handlers
	dg.AddHandler(onReady)
	dg.AddHandler(onVoiceStateUpdate)
	dg.AddHandler(onGuildScheduledEventCreate)
	dg.AddHandler(onGuildScheduledEventUpdate)
	dg.AddHandler(onGuildScheduledEventDelete)
	dg.AddHandler(onGuildScheduledUserAdd)
	dg.AddHandler(onGuildScheduledUserRemove)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}
	defer dg.Close()

	log.Println("Bot is now running. Press CTRL-C to exit.")
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}
