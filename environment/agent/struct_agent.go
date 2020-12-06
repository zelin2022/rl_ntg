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

func contains(list []int, item int)(bool){
  for i := range list {
    if list[i] == item{
      return true
    }
  }
  return false
}
func DeleteAgents(agents []Agent, toDelete []int)[]Agent{
  // logic here is confusing...
  // but basically: we are trying to overwrite to-be-removed positions with useful data from the end
  // j represents position of data from the end
  // i belong to toDelete, toDelete represetns positions to be removed

  var len_after_delete = len(agents) - len(toDelete)
  for i,j := 0, len_after_delete; i < len(toDelete); i, j = i+1, j+1{
    if contains(toDelete, j){ // if we need to delete J anyway
      // skip but do not increment i
      i--
    }else{
      if (toDelete[i] >= len_after_delete){ // if this one will be deleted, we don't save it
        // skip
      }else{
        agents[toDelete[i]] = agents[j]
      }
    }
  }
  return agents[:len_after_delete]
}

func GetAllAgentIDs(agents []Agent) []string {
  var IDs []string
  for i := range agents{
    IDs = append(IDs, agents[i].ID)
  }
  return IDs
}
