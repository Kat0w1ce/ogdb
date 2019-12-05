package example

import (
	"fmt"
	"github.com/tecbot/gorocksdb"
)

func Run()  int {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, "/path/to/db")
	if err!=nil{
		fmt.Println("open db error")
		return 0
	}
	ro := gorocksdb.NewDefaultReadOptions()
	wo := gorocksdb.NewDefaultWriteOptions()
	// if ro and wo are not used again, be sure to Close them.
	err = db.Put(wo, []byte("foo"), []byte("bar"))
	value, err := db.Get(ro, []byte("foo"))
	defer value.Free()
	err = db.Delete(wo, []byte("foo"))
}