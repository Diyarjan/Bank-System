package listeners

import (
	"encoding/json"
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/pkg/repository"
	"github.com/Diyarjan/BankSystem/third_party/kafkaPart"
	"log"
)

type TransferConsumer struct {
	repo     repository.Transaction
	consumer *kafkaPart.Consumer
}

func NewTransferConsumer(repo repository.Transaction, groupId string, consumer *kafkaPart.Consumer) *TransferConsumer {
	return &TransferConsumer{repo: repo, consumer: consumer}
}

func (d *TransferConsumer) StartListening() {
	var transferStruct BankSystem.Transfer
	for {
		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message - %s\n", err)
			continue
		}

		if err := json.Unmarshal(msg.Value, &transferStruct); err != nil {
			log.Printf("Failed to unmarshal message - %s\n", err)
			continue
		}

		newBalance, err := d.repo.TransferFunds(transferStruct)
		if err != nil {
			log.Printf("Failed to deposit to account - %s\n", err)
			continue
		}

		log.Printf("Transfer to account: %d successfully your new balance %f", transferStruct.ToAccountID, newBalance)
	}
}
