package main

import (
	"bufio"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
)

var (
	token string
	cid   string
)

func init() {
	// Make sure token is prefixed with "Bot "
	token = os.Getenv("DISCORD_TOKEN")
	// ID of the channel you want to clean
	cid = os.Getenv("DISCORD_CHANNEL_ID")
}

func main() {
	counter := 0
	// Setup Discord
	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatal("[ERROR] Couldn't create Discord session,", err)
		return
	}

	// Safety measure to make sure we're purging the right channel
	lastMessages, err := dg.ChannelMessages(cid, 1, "", "")
	if err != nil || len(lastMessages) == 0 {
		log.Fatal("[ERROR] Failed to query latest message:", err)
	}

	fmt.Println("Most recent message in this channel:")
	fmt.Printf("%s: %s\n",
		lastMessages[0].Author.Username, lastMessages[0].Content)

	if !confirm("Delete all messages in this channel?") {
		return
	}

	for {
		// Fetch and delete 100 (the max amount) messages at a time
		msgs, err := dg.ChannelMessages(cid, 100, "", "")
		if err != nil {
			log.Println("[WARN] Couldn't fetch channel messages:", err)
		}

		if len(msgs) == 0 {
			break
		}

		err = dg.ChannelMessagesBulkDelete(cid, extractIDs(msgs))
		if err != nil {
			log.Println("[WARN] Failed to delete messages, but carrying on:", err)
		} else {
			fmt.Printf(".")
			counter += len(msgs)
		}
	}
	fmt.Println()
	fmt.Printf("Deleted %d messages in total!\n", counter)
}

func confirm(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [Y/n] ", s)
	resp, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("[ERROR]", err)
	}

	resp = strings.ToLower(strings.TrimSpace(resp))
	return (resp == "y" || resp == "yes")
}

func extractIDs(msgs []*discordgo.Message) []string {
	ids := make([]string, 0)
	for _, m := range msgs {
		ids = append(ids, m.ID)
	}
	return ids
}
