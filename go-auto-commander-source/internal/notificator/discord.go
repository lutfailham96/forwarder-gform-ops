package notificator

import (
	"github.com/bwmarrin/discordgo"
)

// Discord represents a Discord notificator.
type Discord struct {
	TokenBot    string
	RecepientId string
	Message     string
	Session     *discordgo.Session
}

// newDiscordSession creates a new Discord session.
func newDiscordSession(tokenBot string) (*discordgo.Session, error) {
	return discordgo.New("Bot " + tokenBot)
}

// newDiscordConnection opens a websocket connection to Discord.
func newDiscordConnection(dg *discordgo.Session) error {
	return dg.Open()
}

// getDiscordId gets the Discord ID of a user.
func getDiscordId(dg *discordgo.Session, username string) (string, error) {
	user, err := dg.User("@" + username)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

// NewDiscord creates a new Discord notificator.
func NewDiscord(tokenBot, recepientId string) (*Discord, error) {
	// create a new Discord session
	dg, err := newDiscordSession(tokenBot)
	if err != nil {
		return nil, err
	}

	// get Discord ID of the user
	exactRecepientId, err := getDiscordId(dg, recepientId)
	if err != nil {
		exactRecepientId = recepientId
	}

	return &Discord{
		TokenBot:    tokenBot,
		RecepientId: exactRecepientId,
		Message:     "",
		Session:     dg,
	}, nil
}

// SendMessage sends a message to the user.
func (d *Discord) SendMessage(message string) error {
	// Open a websocket connection to Discord.
	err := newDiscordConnection(d.Session)
	if err != nil {
		return err
	}

	// Get the user's DM channel.
	channel, err := d.Session.UserChannelCreate(d.RecepientId)
	if err != nil {
		return err
	}

	// Send the message to the user.
	d.Message = message
	_, err = d.Session.ChannelMessageSend(channel.ID, d.Message)

	return err
}
