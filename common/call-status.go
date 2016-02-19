package common

import(
    "encoding/json"
)

type CallStatus struct {
    To string
    CallStatus string
    CallSid string
}

func NewCallStatus(jsonStr string) *CallStatus {
    var cs CallStatus
    json.Unmarshal([]byte(jsonStr), &cs)
    return &cs
}

func (self *CallStatus) ToJson() string {
    bytes, _ := json.Marshal(&self)
    return string(bytes)
}