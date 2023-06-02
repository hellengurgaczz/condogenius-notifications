package main

import (
	"Condogenius-notifications/db"
	"Condogenius-notifications/models"
	"Condogenius-notifications/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Conexão com banco de dados
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Houve algum problema na conexão com banco de dados")
	} else {
		fmt.Println("Conexão com banco de dados estabelecida")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "handy-courage-388421")
	if err != nil {
		log.Fatalf("Falha ao criar o cliente do Pub/Sub: %v", err)
	}
	subscription := client.Subscription("send-notifications-sub")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-interrupt:
			fmt.Println("Programa encerrado.")
			return
		default:
			err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
				fmt.Printf("Mensagem Recebida: %s\n", string(msg.Data))

				var notification models.Notification
				err := json.Unmarshal(msg.Data, &notification)
				if err != nil {
					log.Printf("Erro ao converter a mensagem em Notification: %v", err)
					msg.Nack()
					return
				}

				err = repository.Save(db, notification)
				if err != nil {
					log.Printf("Erro ao salvar a mensagem no banco de dados: %v", err)
					msg.Nack()
					return
				}

				msg.Ack()
			})
			if err != nil {
				log.Printf("Erro ao receber mensagens: %v", err)
				time.Sleep(5 * time.Second) // Aguardar antes de tentar novamente em caso de erro
			}
		}
	}
}
