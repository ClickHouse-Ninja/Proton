package server

import (
	"log"
	"time"
)

type dict struct {
	id            uint64
	column, value string
}

func (server *server) backgroundDictionary() {
	for tick := time.Tick(time.Second); ; {
		block := server.dictBlock.Copy()
	loop:
		for block.Reserve(); ; {
			select {
			case dict := <-server.dictBacklog:
				block.NumRows++
				block.WriteUInt64(0, dict.id)
				block.WriteString(1, dict.value)
				block.WriteString(2, dict.column)
			case <-tick:
				break loop
			}
		}
		if err := server.writeBlock(insertIntoDictionarySQL, block); err != nil {
			log.Println("dictionary write error: ", err)
		}
	}
}

const (
	insertIntoDictionarySQL = `
	INSERT INTO proton.dictionary (
		ID
		, Value
		, Column
	) VALUES (
		?, ?, ?
	)
	`
)
