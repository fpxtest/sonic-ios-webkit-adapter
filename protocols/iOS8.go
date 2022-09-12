package protocols

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
	adapters "sonic-ios-webkit-adapter/adapter"
)

type iOS8 struct {
	adapter *adapters.Adapter
}

func initIOS8(protocol *ProtocolAdapter) {
	result := &iOS8{
		adapter: protocol.adapter,
	}
	protocol.init()
	protocol.adapter.AddMessageFilter("error", result.targetError)
	protocol.mapSelectorList = result.mapSelectorList
}

func (i *iOS8) targetError(message []byte) []byte {
	params := map[string]interface{}{
		"id":     gjson.Get(string(message), "id"),
		"result": map[string]interface{}{},
	}
	msg, err := sjson.Set("", "", params)
	if err != nil {
		log.Panic(err)
	}
	return []byte(msg)
}

func (i *iOS8) mapSelectorList(selectorList gjson.Result, message string) string {
	cssRange := selectorList.Get("range")
	var err error
	for _, selector := range selectorList.Get("selectors").Array() {
		message, err = sjson.Set(message, selector.Path(message), map[string]interface{}{
			"text": selector.Value(),
		})
		if cssRange.Exists() {
			message, err = sjson.Set(message, selector.Get("range").Path(message), cssRange.Value())
		}
		if err != nil {
			log.Panic(err)
		}
	}
	message, err = sjson.Delete(message, cssRange.Path(message))
	if err != nil {
		log.Panic(err)
	}
	return message
}
