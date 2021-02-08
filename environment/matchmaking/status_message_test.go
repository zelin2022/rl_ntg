package matchmaking
import (
  "testing"
)


func TestToStatusBody(t *testing.T){
  str := `{"queue": "amq.gen-N_v5tBrTqYsnUYClAvIIYw", "mmc": "fakeMMC"}`
  output, err := toStatusBody(str)
  if err != nil ||  output.Queue != "amq.gen-N_v5tBrTqYsnUYClAvIIYw" || output.MMC != "fakeMMC" {
    t.Error("toStatusBody failed to parse json string into struct")
  }
}
