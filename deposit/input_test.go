package deposit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/testutil"
	"github.com/umbracle/ethgo/wallet"
	"github.com/umbracle/go-eth-consensus/bls"
)

func TestDeposit_Signing(t *testing.T) {
	kk := bls.NewRandomKey()
	data, err := Input(kk, nil, ethgo.Gwei(MinGweiAmount).Uint64())
	if err != nil {
		t.Fatal(err)
	}

	err = Verify(data)
	require.NoError(t, err)
}

func TestDeposit_EndToEnd(t *testing.T) {
	server := testutil.NewTestServer(t, nil)
	defer server.Close()

	ecdsaKey, _ := wallet.GenerateKey()
	server.Transfer(ecdsaKey.Address(), ethgo.Ether(MinGweiAmount+1))

	// deploy the contract
	receipt, err := server.SendTxn(&ethgo.Transaction{
		Input: DepositBin(),
	})
	assert.NoError(t, err)

	client, _ := jsonrpc.NewClient(server.HTTPAddr())
	code, err := client.Eth().GetCode(receipt.ContractAddress, ethgo.Latest)
	assert.NoError(t, err)
	assert.NotEqual(t, code, "0x")

	// sign the deposit
	key := bls.NewRandomKey()

	input, err := Input(key, nil, ethgo.Gwei(MinGweiAmount).Uint64())
	assert.NoError(t, err)

	// deploy transaction
	depositContract := NewDeposit(receipt.ContractAddress, contract.WithSender(ecdsaKey), contract.WithJsonRPC(client.Eth()))

	txn, err := depositContract.Deposit(input.Pubkey[:], input.WithdrawalCredentials[:], input.Signature[:], input.Root)
	assert.NoError(t, err)

	txn.WithOpts(&contract.TxnOpts{Value: ethgo.Ether(MinGweiAmount)})

	assert.NoError(t, txn.Do())

	_, err = txn.Wait()
	assert.NoError(t, err)

	// query the contract
	count, err := depositContract.GetDepositCount()
	assert.NoError(t, err)
	assert.Equal(t, int(count[0]), 1)
}
