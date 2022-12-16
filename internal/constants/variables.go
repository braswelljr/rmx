package constants

import (
	"math/rand"
	"time"
)

var (
	DotChar         = 46                                                                                  // DotChar - the ASCII character for a dot
	SeedSrc         = rand.NewSource(time.Now().UnixNano())                                               // random seed
	SourceKey       = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+{}|:<>?" // the key to be used for the generation
	LetterIndexBits = 6                                                                                   // 6 bits to represent a letter index
	LetterIndexMask = 1<<LetterIndexBits - 1                                                              // All 1-bits, as many as letterIndexBits
)
