package config

var BlockPipelineWidth = 31

var IsFastPeer = true
var IsStorage = false
var IsEndorser = false
var IsBenchmark = false

var RegisterBlockStore func(ledgerId string, blockStore interface{})
