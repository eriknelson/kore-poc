package comm

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// NOTE: This is a mock client that mocks out a vendored API that serves as
// a client into a particular platform, say Discord or IRC. Normally each
// platform's client would look completely different, but this is an Stdin
// based client that demos dynamic, incomming content.

type MockChatMessage struct {
	User    string
	Message string
}

type MockPlatformClient struct {
	Chat chan MockChatMessage

	name string
}

func NewMockPlatformClient(name string) *MockPlatformClient {
	return &MockPlatformClient{
		name: name,
		Chat: make(chan MockChatMessage),
	}
}

func (c *MockPlatformClient) Connect() {
	go func() {
		log.Debug("MockPlatformClient::Connect")

		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			platformName, chatMsg, err := c.structuredMsg(text)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			log.Infof("Got stdin msg bound for adapter: %s", platformName)
			if platformName == c.name {
				c.Chat <- *chatMsg
			} else {
				log.Infof("%s client ignoring stdin msg...", c.name)
			}
		}
	}()
}

func (c *MockPlatformClient) structuredMsg(text string) (string, *MockChatMessage, error) {
	split := strings.Split(text, " ")
	if !(len(split) > 1) {
		return "", nil, errors.New("Must send stdin message in format of '<adapter_name> <content>'")
	}

	adapter := split[0]
	strippedSplit := split[1:len(split)]
	message := strings.Join(strippedSplit, " ")

	return adapter, &MockChatMessage{
		User:    fmt.Sprintf("%s-user", c.name),
		Message: message,
	}, nil
}

func (c *MockPlatformClient) SendMessage(m string) {
	log.Infof("Discord client got a message from the adapter! [ %s ]", m)
}
