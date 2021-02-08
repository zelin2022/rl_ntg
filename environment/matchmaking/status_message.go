package matchmaking

import(
  "encoding/json"
)

type StatusMessageBody struct{
  Queue string `json:"queue"`
  MMC string `json:"mmc"`
}
func toStatusBody(jsonString string)(StatusMessageBody, error){
  var outputStruct StatusMessageBody
  err := json.Unmarshal([]byte(jsonString), &outputStruct)
  return outputStruct, err
}
