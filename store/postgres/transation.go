package postgres

import (
	"encoding/binary"
	"time"
	"unsafe"

	"github.com/d7561985/mongo-ab/pkg/changing"
)

const sizeUint64 = int(unsafe.Sizeof(uint64(0)))

type Balance struct {
	AccountID     uint64
	Balance       float64
	DepositAllSum float64
	DepositCount  int32

	PincoinBalance float64
	PincoinsAllSum float64
}

type Journal struct {
	// for csv use hexadecimal
	ID []byte `csv:"id"`
	// for csv use hexadecimal via s3 load
	TransactionID []byte `csv:"transactionId"`

	AccountID uint64 `csv:"accountId"`

	CreatedAt time.Time `csv:"created_at"`

	Balance float64 `csv:"balance"`
	Change  float32 `csv:"change"`

	PincoinBalance float64 `csv:"pincoinBalance"`
	PincoinChange  float32 `csv:"pincoinChange"`

	Type OpType `csv:"type"`

	Project  Project `csv:"project"`
	Currency int8    `csv:"currency"`

	Revert bool `csv:"revert"`
}

func NewJournal(b Balance, in changing.Transaction) Journal {
	trn, err := GetTransactionID(in)
	if err != nil {
		trn = []byte{}
	}

	return Journal{
		ID:            in.Set.ID[:],
		TransactionID: trn,
		AccountID:     b.AccountID,

		CreatedAt: time.Now(),

		Balance:        b.Balance,
		Change:         float32(in.Change),
		PincoinBalance: b.PincoinBalance,
		PincoinChange:  float32(in.PincoinChange),

		Type:     NewOperationType(in.Type),
		Project:  NewProject(in.Project),
		Currency: int8(in.Currency),
		Revert:   false,
	}
}

func GetTransactionID(r changing.Transaction) ([]byte, error) {
	if r.TransactionID > 0 {
		tID := make([]byte, sizeUint64)
		binary.LittleEndian.PutUint64(tID, r.TransactionID)

		return tID, nil
	}

	return r.TransactionIDBson[:], nil
}
