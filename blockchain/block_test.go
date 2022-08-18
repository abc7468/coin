package blockchain

import (
	"coin/utils"
	"reflect"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				return nil
			},
		}
		_, err := FindBlock("xx")
		if err == nil {
			t.Error("The block should not be found!")
		}
	})
	t.Run("Block is found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{}
				return utils.ToBytes(b)
			},
		}
		b, err := FindBlock("xx")
		if err != nil {
			t.Error("The block should be found")
		}
		if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
			t.Error("The block should be found")
		}
	})
}
