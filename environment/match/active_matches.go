package match

import (
  "sync"
)

type ActiveMatches struct{
  Matches []Match
  Mutex sync.Mutex
}
