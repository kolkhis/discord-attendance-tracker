package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	dg.AddHandler(onGuildScheduledEventUserAdd)
	dg.AddHandler(onGuildScheduledEventUserRemove)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}
	defer dg.Close()

	log.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

func onVoiceStateUpdate(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	log.Printf("Voice state updated for user %v in guild %v", event.UserID, event.GuildID)
}

func onGuildScheduledEventCreate(s *discordgo.Session, event *discordgo.GuildScheduledEventCreate) {
	log.Printf("Scheduled event created: %v in guild %v", event.GuildScheduledEvent.Name, event.GuildID)
}

func onGuildScheduledEventUpdate(s *discordgo.Session, event *discordgo.GuildScheduledEventUpdate) {
	log.Printf("Scheduled event updated: %v in guild %v", event.GuildScheduledEvent.Name, event.GuildID)
}

func onGuildScheduledEventDelete(s *discordgo.Session, event *discordgo.GuildScheduledEventDelete) {
	log.Printf("Scheduled event deleted: %v in guild %v", event.GuildScheduledEvent.Name, event.GuildID)
}

func onGuildScheduledEventUserAdd(s *discordgo.Session, event *discordgo.GuildScheduledEventUserAdd) {
	log.Printf("User %v added to scheduled event %v in guild %v", event.UserID, event.GuildScheduledEventID, event.GuildID)
}

func onGuildScheduledEventUserRemove(s *discordgo.Session, event *discordgo.GuildScheduledEventUserRemove) {
	log.Printf("User %v removed from scheduled event %v in guild %v", event.UserID, event.GuildScheduledEventID, event.GuildID)
}
