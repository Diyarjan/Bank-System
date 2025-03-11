package listeners

import (
	"encoding/json"
	"github.com/Diyarjan/BankSystem"
	"github.com/Diyarjan/BankSystem/pkg/repository"
	"github.com/Diyarjan/BankSystem/third_party/kafkaPart"
	"log"
)

type WithdrawConsumer struct {
	repo     repository.Transaction
	consumer *kafkaPart.Consumer
}

func NewWithdrawConsumer(repo repository.Transaction, groupId string, consumer *kafkaPart.Consumer) *WithdrawConsumer {
	return &WithdrawConsumer{repo: repo, consumer: consumer}
}

func (d *WithdrawConsumer) StartListening() {
	var withdrawStruct BankSystem.DebitCreditStruct
	for {
		msg, err := d.consumer.PollMessage()
		if err != nil {
			log.Printf("Failed to read message - %s\n", err)
			continue
		}
		if err := json.Unmarshal(msg.Value, &withdrawStruct); err != nil {
			log.Printf("Failed to unmarshal message - %s\n", err)
			continue
		}

		err = d.repo.WithdrawFromAccount(withdrawStruct)
		if err != nil {
			log.Printf("Failed to withdraw from account - %s\n", err)
		}
		log.Printf("Withdraw from account: %d successfully", withdrawStruct.AccountID)
	}

}
