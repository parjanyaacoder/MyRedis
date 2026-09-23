package config

var Host string = "0.0.0.0"

var Port int = 2610
var KeysLimit int = 100

var EvictionRatio float64 = 0.40

var AOFFile string = "./MyRedis.aof"

var EvictionStrategy string = "allkeys-random"
