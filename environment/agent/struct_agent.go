package agent
import (
  "time"
)

type Agent struct {
  ID string
  Queue string
  LastActive time.Time
}


func (a *Agent)RenewActive(){
  a.LastActive = time.Now()
}

// HELPER METHODS ==============================================================================
func FindAgent(agents []Agent, agentID string) (bool, int) {
  for i := 0; i < len(agents); i++{
    if agentID == agents[i].ID {
      return true, i
    }
  }
  return false, 0
}

func DeleteAgent(agents []Agent, position int) []Agent {
  agents[position] = agents[len(agents) - 1] // swap last element to element to delete
  return agents[:len(agents) - 1] // return slice with last element excluded
}
