package log

import "encoding/json"

//zap中使用zap.Any打印数据之前将其转换为原始json
func AnyObject(obj interface{}) *json.RawMessage {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	p := json.RawMessage(bytes)
	return &p
}
