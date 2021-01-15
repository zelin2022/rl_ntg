package agent
import (
  "testing"
)
func TestFindAgent(t *testing.T){
  // populate
  var agents []Agent
  agents = append(agents, Agent{
    ID: "testID0",
  })
  agents = append(agents, Agent{
    ID: "testID1",
  })
  // agent in slice be found
  out_0_0, out_0_1 := FindAgent(agents, "testID1")
  if out_0_0 != true || out_0_1 != 1 {
    t.Error("FindAgent failed to find an exiting agent")
  }
  // agent not in slice not found
  out_1_0, out_1_1 := FindAgent(agents, "testID2")
  if out_1_0 != false || out_1_1 != -1 {
    t.Error("FindAgent failed to not find a non-exiting agent")
  }
}

func TestDeleteAgent(t *testing.T){
  // populate
  var agents []Agent
  agents = append(agents, Agent{
    ID: "testDelete_0",
  })
  // agent in slice is deleted
  var err error
  agents, err = DeleteAgent(agents, 0)
  if len(agents) != 0 || err != nil {
    t.Error("DeleteAgent failed to delete existing agent")
  }
  agents, err = DeleteAgent(agents, 0)
  if len(agents) != 0 || err == nil {
    t.Error("DeleteAgent failed to return error on out of range agent")
  }
}

func TestDeleteAgentByID(t *testing.T){
  // populate
  var agents []Agent
  agents = append(agents, Agent{
    ID: "testDelete_0",
  })
  agents = append(agents, Agent{
    ID: "testDelete_1",
  })
  // agent in slice is deleted
  var err error
  agents, err = DeleteAgentByID(agents, "testDelete_0")
  if len(agents) != 1 || err != nil || agents[0].ID != "testDelete_1" {
    t.Error("DeleteAgentByID failed to delete an existing agent")
  }
  // ID not exist returns error
  agents, err = DeleteAgentByID(agents, "testDelete_2")
  if len(agents) != 1 || err == nil || agents[0].ID != "testDelete_1" {
    t.Error("DeleteAgentByID failed to fail when deleting agent by a non-existing ID")
  }
}

func Testcontains(t *testing.T){
  var mylist []int
  mylist = append(mylist, 3)
  mylist = append(mylist, 444)
  mylist = append(mylist, 82382)
  var ret bool
  // returns true if list contains int
  ret = contains(mylist, 444)
  if ret != true {
    t.Error("contains failed to return true when list contains int")
  }
  ret = contains(mylist, 4)
  if ret != false {
    t.Error("contains failed to return false when list does not contain int")
  }
}

func TestDeleteAgents(t *testing.T){
  // populate
  var agents []Agent
  agents = append(agents, Agent{
    ID: "testDelete_0",
  })
  agents = append(agents, Agent{
    ID: "testDelete_1",
  })
  agents = append(agents, Agent{
    ID: "testDelete_2",
  })
  var err error
  // deletes multiple agents when I pass it multiple positions
  toDelete := []int{0, 2}
  agents, err = DeleteAgents(agents, toDelete)
  if len(agents) != 1 || agents[0].ID != "testDelete_1" || err != nil {
    t.Error("DeleteAgents failed to delete agents correctly")
  }
  // returns an error when positions passed to it is invalid
  agents, err = DeleteAgents(agents, toDelete)
  if len(agents) != 1 || agents[0].ID != "testDelete_1" || err == nil {
    t.Error("DeleteAgents failed to return error when invalid positions are passed to it")
  }
}

func TestGetAllAgentIDs(t *testing.T) {
  // populate
  var agents []Agent
  agents = append(agents, Agent{
    ID: "q1!",
  })
  agents = append(agents, Agent{
    ID: "@2w",
  })
  agents = append(agents, Agent{
    ID: "aJd#2jkdio.vsd",
  })
  // returns slice of string containing correct ids
  ret0 := GetAllAgentIDs(agents)
  if len(ret0) != 3 || ret0[0] != "q1!" || ret0[1] != "@2w" || ret0[2] != "aJd#2jkdio.vsd" {
    t.Error("GetAllAgentIDs failed to return the correct string slice")
  }
  // return empty slice when passing empty slice
  agents = []Agent{}
  ret1 := GetAllAgentIDs(agents)
  if len(ret1) != 0 {
    t.Error("GetAllAgentIDs failed to return empty slice when passing empty slice")
  }
}
