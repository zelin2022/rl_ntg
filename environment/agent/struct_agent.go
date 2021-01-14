package agent
import (
  "errors"
)

type Agent struct {
  ID string
  Queue string
}

// HELPER METHODS ==============================================================================

// finds first currance by matching ID
func FindAgent(agents []Agent, agentID string) (bool, int) {
  for i := range agents{
    if agentID == agents[i].ID {
      return true, i
    }
  }
  return false, 0
}

// delete agent at position, return deleted slice
func DeleteAgent(agents []Agent, position int) []Agent {
  agents[position] = agents[len(agents) - 1] // swap last element to element to delete
  return agents[:len(agents) - 1] // return slice with last element excluded
}

func DeleteAgentByID(agents []Agent, id string)([]Agent, error){
  for i := range agents{
    if agents[i].ID == id{
      agents[i] = agents[len(agents)-1]
      return agents[:len(agents)-1], nil
    }
  }
  return agents, errors.New("ID not found, ID: " + id)
}

// check if a slice of ints contains a specific int
func contains(list []int, item int)(bool){
  for i := range list {
    if list[i] == item{
      return true
    }
  }
  return false
}

// delete a subslice of agents in a slice, returns the modified slice
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

// return a slice of string which are agent IDs mapping the slice of agents
func GetAllAgentIDs(agents []Agent) []string {
  var IDs []string
  for i := range agents{
    IDs = append(IDs, agents[i].ID)
  }
  return IDs
}
