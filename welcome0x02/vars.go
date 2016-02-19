package main

import (
	"math/rand"
	"time"
)

var randomSource = rand.NewSource(time.Now().UnixNano())
var randoms = rand.New(randomSource)
