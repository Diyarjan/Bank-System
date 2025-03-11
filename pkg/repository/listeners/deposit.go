package listeners

import (
	"encoding/json"
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/pkg/repository"
	"github.com/Diyarjan/BankSystem/third_party/kafkaPart"
	"log"
)

type DepositConsumer struct {
	repo     repository.Transaction
	consumer *kafkaPart.Consumer
}

func NewDepositConsumer(repo repository.Transaction, groupId string, consumer *kafkaPart.Consumer) *DepositConsumer {
	return &DepositConsumer{repo: repo, consumer: consumer}
}

func (d *DepositConsumer) StartListening() {
	var depositStruct BankSystem.DebitCreditStruct
	for {

		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message - %s\n", err)
			continue
		}

		if err := json.Unmarshal(msg.Value, &depositStruct); err != nil {
			log.Printf("Failed to unmarshal message - %s\n", err)
			continue
		}

		if err := d.repo.DepositToAccount(depositStruct); err != nil {
			log.Printf("Failed to deposit to account - %s\n", err)
		}

		log.Printf("Deposit to account: %d successfully yee", depositStruct.AccountID)

	}

}
