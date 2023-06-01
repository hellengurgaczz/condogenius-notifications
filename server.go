package main

import (
	"context"
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
	// Configurar a conexão com o banco de dados MySQL
	db, err := db.connectDB()
	defer db.Close()

	ctx := context.Background()

	// Crie um cliente do Pub/Sub
	client, err := pubsub.NewClient(ctx, "handy-courage-388421")
	if err != nil {
		log.Fatalf("Falha ao criar o cliente do Pub/Sub: %v", err)
	}

	// Crie uma assinatura para receber mensagens
	subscription := client.Subscription("send-notifications-sub")

	// Crie um canal para capturar sinais de interrupção
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Execute o loop para receber mensagens continuamente
	for {
		select {
		case <-interrupt:
			// Encerre o programa caso receba um sinal de interrupção
			fmt.Println("Programa encerrado.")
			return
		default:
			err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
				// Processar a mensagem recebida
				fmt.Printf("Mensagem Recebida: %s\n", string(msg.Data))

				// Confirme o recebimento da mensagem (marcando-a como concluída)
				msg.Ack()
			})
			if err != nil {
				log.Printf("Erro ao receber mensagens: %v", err)
				time.Sleep(5 * time.Second) // Aguardar antes de tentar novamente em caso de erro
			}
		}
	}
}
