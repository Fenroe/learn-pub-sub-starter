package main

import (
	"fmt"
	"log"

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
	fmt.Println("Starting Peril client...")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("Could not get username: %v\n", err)
	}
	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		1,
	)
	if err != nil {
		log.Printf("Could not bind queue to message exchange: %v\n", err)
	}

	gameState := gamelogic.NewGameState(username)
	programIsRunning := true
	for programIsRunning {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			fmt.Println(input[1:][1])
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Printf("Could not spawn unit: %v\n", err)
				continue
			}
		case "move":
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("Could not move unit: %v\n", err)
				continue
			}
			fmt.Println("Moved piece")
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			programIsRunning = false
		default:
			fmt.Println("Could not execute command")
			fmt.Println("Enter 'help' for a list of available commands")
		}
	}
}
