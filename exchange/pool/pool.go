package factory

import (
	"sync"

	"github.com/opefago/bot-tave/exchange/base"
)

type pool struct {
	idle   []base.Exchange
	active []base.BaseExchange

	capacity int
	mutex    sync.RWMutex
}

func Init(path string) *pool {
	return &pool{}
}
