package btypes

type TradeDirection string

const (
	TradeDirectionCall TradeDirection = "call"
	TradeDirectionPut  TradeDirection = "put"
)

type TradeShouldWaitForResult bool

const (
	WaitForResult      TradeShouldWaitForResult = true
	DoNotWaitForResult TradeShouldWaitForResult = false
)
