package blockchain

import (
	"coin/utils"
	"reflect"
	"sync"
	"testing"
)

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}
func (fakeDB) SaveBlock(hash string, data []byte) {
	return
}
func (fakeDB) SaveChain(data []byte) {
	return
}
func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}
func (fakeDB) DeleteAllBlocks() {
	return
}

func TestBlockchain(t *testing.T) {
	t.Run("불러올 블록체인 없을 때", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				return nil
			},
		}
		bc := Blockchain()
		if bc.Height != 1 {
			t.Error("Height가 1인 블록체인이 만들어 져야합니다.")
		}
	})
	t.Run("불러올 블록체인 있을 때", func(t *testing.T) {
		once = sync.Once{}
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				b := &blockchain{
					Height: 2,
				}
				return utils.ToBytes(b)
			},
		}
		bc := Blockchain()
		if bc.Height != 2 {
			t.Errorf("Height가 %d인 블록체인이 불려져야 하지만, Height가 %d인 블록체인이 불려집니다.", 2, bc.Height)
		}
		if reflect.TypeOf(bc) != reflect.TypeOf(&blockchain{}) {
			t.Error("Blockchain()이 반환해야 하는 타입이 옳지 않습니다.")
		}
	})
}
