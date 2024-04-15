package brokerws

import (
	"fmt"
	"time"

	"github.com/patrickkdev/Go-IQOption-API/btypes"
	"github.com/patrickkdev/Go-IQOption-API/internal/debug"
	"github.com/patrickkdev/Go-IQOption-API/internal/tjson"
)

type binaryTradeData struct {
	Name             string              `json:"name"`
	MicroserviceName string              `json:"microserviceName"`
	Msg              btypes.BinaryResult `json:"msg"`
}

type digitalTradeData struct {
	Name             string               `json:"name"`
	MicroserviceName string               `json:"microserviceName"`
	Msg              btypes.DigitalResult `json:"msg"`
}

func (ws *Socket) CheckResultBinary(tradeID int, timeout time.Time) (*btypes.BinaryResult, bool, error) {
	debug.IfVerbose.Printf("Calling check result binary with tradeID: %d\n", tradeID)

	var err error = nil
	var res binaryTradeData

	ws.AddEventListener("position-changed", func(event []byte) {
		res, err = tjson.Unmarshal[binaryTradeData](event)

		debug.IfVerbose.Println("Position changed for binary trade")
		debug.IfVerbose.PrintAsJSON(res)
	})

	for res.Msg.ExternalID != tradeID || res.Msg.Status != "closed" {
		time.Sleep(time.Second)

		if ws.Closed {
			err = fmt.Errorf("websocket closed")
			break
		}

		if time.Since(timeout) > 0 {
			err = fmt.Errorf("timed out waiting for response")
			break
		}

		debug.IfVerbose.Println("Position changed for binary trade:", res.Msg.Status, res.Msg.ExternalID, tradeID)
	}

	ws.RemoveEventListener("position-changed")

	win := res.Msg.CloseReason == "win"
	debug.IfVerbose.Println("Result:", res.Msg.Status, res.Msg.ExternalID, tradeID, win)

	if err != nil {
		return nil, false, err
	}

	return &res.Msg, res.Msg.CloseReason == "win", err
}

func (ws *Socket) CheckResultDigital(tradeID int, timeout time.Time) (*btypes.DigitalResult, bool, error) {
	debug.IfVerbose.Printf("Calling check result digital with tradeID: %d\n", tradeID)

	var err error = nil
	var res digitalTradeData

	ws.AddEventListener("position-changed", func(event []byte) {
		res, err = tjson.Unmarshal[digitalTradeData](event)

		debug.IfVerbose.Println("Position changed for digital trade")
		debug.IfVerbose.PrintAsJSON(res)
	})

	for {
		time.Sleep(time.Second)

		if ws.Closed {
			err = fmt.Errorf("websocket closed")
			break
		}

		if time.Since(timeout) > 0 {
			err = fmt.Errorf("timed out waiting for response")
		}

		orderIDs := res.Msg.RawEvent.DigitalOptionsPositionChanged1.OrderIds

		if len(orderIDs) == 0 {
			continue
		}

		orderID := res.Msg.RawEvent.DigitalOptionsPositionChanged1.OrderIds[0]

		if orderID != tradeID {
			continue
		}

		debug.IfVerbose.Println("Position changed for digital trade:", res.Msg.Status, orderID, tradeID)

		if res.Msg.Status != "closed" {
			continue
		}

		break
	}

	ws.RemoveEventListener("position-changed")

	win := res.Msg.Pnl > 0
	debug.IfVerbose.Println("Check result digital: ", res.Msg.Pnl, win)

	if err != nil {
		return nil, false, err
	}

	return &res.Msg, win, err
}
