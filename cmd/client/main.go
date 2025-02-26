package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connectionString := "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("", err)
	}
	defer connection.Close()

	connectionChannel, err := connection.Channel()
	if err != nil {
		log.Fatal("", err)
	}

	fmt.Println("Starting Peril client...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("Could not get username: %v\n", err)
	}

	pauseUsernameQueueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)

	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		pauseUsernameQueueName,
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		log.Printf("Could not bind queue to message exchange: %v\n", err)
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		pauseUsernameQueueName,
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gameState),
	)

	if err != nil {
		log.Println("Could not establish connection to server: ", err)
	}

	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		"army_moves."+username,
		"army_moves.*",
		pubsub.Transient,
		handlerMove(gameState, connectionChannel),
	)

	if err != nil {
		log.Println("Could not establish connection to server: ", err)
	}

	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		"war",
		"war.*",
		pubsub.Durable,
		handlerWar(gameState, connectionChannel),
	)

	if err != nil {
		log.Println("Could not establish connection to server: ", err)
	}

	programIsRunning := true
	for programIsRunning {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Printf("Could not spawn unit: %v\n", err)
				continue
			}
		case "move":
			move, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("Could not move unit: %v\n", err)
				continue
			}
			pubsub.PublishJSON(connectionChannel, routing.ExchangePerilTopic, "army_moves."+username, move)
			fmt.Println("move published to exchange")
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(input) == 1 {
				fmt.Println("'spam' command needs more words")
				continue
			}
			parsed, err := strconv.ParseInt(input[1], 10, 0)
			if err != nil {
				fmt.Println("'spam' command needs a valid integer")
				continue
			}
			fmt.Printf("Attempting to spam %d messages\n", parsed)
			for i := 0; i < int(parsed); i++ {
				mallog := gamelogic.GetMaliciousLog()
				routingKey := fmt.Sprintf("%s.%s", routing.GameLogSlug, username)
				fmt.Println(routingKey)
				err := pubsub.PublishGob(connectionChannel, routing.ExchangePerilTopic, routingKey, mallog)
				if err != nil {
					fmt.Println(err)
				}
			}
			fmt.Println("Finished spam attempt")
		case "quit":
			gamelogic.PrintQuit()
			programIsRunning = false
		default:
			fmt.Println("Could not execute command")
			fmt.Println("Enter 'help' for a list of available commands")
		}
	}
}
