package matchmaking

import (

)



func Create(channels ChannelBundle){
  var mm MM
  mm.Channels = channels
  go mm.run()
}
