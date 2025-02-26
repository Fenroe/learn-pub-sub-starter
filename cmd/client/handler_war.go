package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

func getGameLogMessage(winner, loser string, outcome gamelogic.WarOutcome) string {
	if outcome == gamelogic.WarOutcomeDraw {
		return fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
	}
	if outcome == gamelogic.WarOutcomeYouWon || outcome == gamelogic.WarOutcomeOpponentWon {
		return fmt.Sprintf("%s won a war against %s", winner, loser)
	}
	return ""
}

func handlerWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {

	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {

		defer fmt.Print("> ")

		outcome, winner, loser := gs.HandleWar(rw)

		routingKey := fmt.Sprintf("%s.%s", routing.GameLogSlug, gs.GetUsername())

		if err := pubsub.PublishGob(
			ch,
			routing.ExchangePerilTopic,
			routingKey,
			routing.GameLog{
				CurrentTime: time.Now().UTC(),
				Message:     getGameLogMessage(winner, loser, outcome),
				Username:    gs.GetUsername(),
			},
		); err != nil {
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
