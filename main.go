package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	guildID := os.Getenv("GUILD_ID")
	outputFile := "members.csv"

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Enable privileged intent to fetch members
	dg.Identify.Intents = discordgo.IntentsGuildMembers

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord connection: %v", err)
	}
	defer dg.Close()

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV headers
	writer.Write([]string{"Username", "Discriminator", "UserID"})

	fmt.Println("Fetching and writing members to CSV...")

	var after string
	limit := 1000

	for {
		members, err := dg.GuildMembers(guildID, after, limit)
		if err != nil {
			log.Fatalf("Error fetching members: %v", err)
		}
		if len(members) == 0 {
			break
		}

		for _, member := range members {
			user := member.User
			record := []string{user.Username, user.Discriminator, user.ID}
			writer.Write(record)
		}

		after = members[len(members)-1].User.ID
	}

	fmt.Printf("Saved all members to %s\n", outputFile)

	// Optional: keep running until Ctrl+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit...")
	<-stop
}
