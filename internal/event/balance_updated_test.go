package event

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mensagem exatamente como o ms-wallet publica no tópico "balances".
const walletBalanceUpdatedMessage = `{"Name":"BalanceUpdated","Payload":{"account_id_from":"from","account_id_to":"to","balance_account_id_from":900,"balance_account_id_to":1100}}`

func TestBalanceUpdated_DecodesWalletMessage(t *testing.T) {
	balanceUpdated := NewBalanceUpdated()

	err := json.Unmarshal([]byte(walletBalanceUpdatedMessage), balanceUpdated)
	assert.NoError(t, err)
	assert.Equal(t, "BalanceUpdated", balanceUpdated.GetName())

	payload, ok := balanceUpdated.GetPayload().(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "from", payload["account_id_from"])
	assert.Equal(t, float64(900), payload["balance_account_id_from"])
}
