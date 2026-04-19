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

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v (User ID: %v)", s.State.User.Username, s.State.User.Discriminator, r.User.ID)
}

func onVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	log.Printf("Voice state updated for user %v (ID: %v) in guild %v", v.Member.User.Username, v.UserID, v.GuildID)
}

func onGuildScheduledEventCreate(s *discordgo.Session, e *discordgo.GuildScheduledEventCreate) {
	log.Printf("Scheduled event created: %v in guild %v", e.GuildScheduledEvent.Name, e.GuildID)
}

func onGuildScheduledEventUpdate(s *discordgo.Session, e *discordgo.GuildScheduledEventUpdate) {
	log.Printf("Scheduled event updated: %v in guild %v", e.GuildScheduledEvent.Name, e.GuildID)
}

func onGuildScheduledEventDelete(s *discordgo.Session, e *discordgo.GuildScheduledEventDelete) {
	log.Printf("Scheduled event deleted: %v in guild %v", e.GuildScheduledEvent.Name, e.GuildID)
}

func onGuildScheduledEventUserAdd(s *discordgo.Session, e *discordgo.GuildScheduledEventUserAdd) {
	log.Printf("User %v added to scheduled event %v in guild %v", e.UserID, e.GuildScheduledEventID, e.GuildID)
}

func onGuildScheduledEventUserRemove(s *discordgo.Session, e *discordgo.GuildScheduledEventUserRemove) {
	log.Printf("User %v removed from scheduled event %v in guild %v", e.UserID, e.GuildScheduledEventID, e.GuildID)
}
