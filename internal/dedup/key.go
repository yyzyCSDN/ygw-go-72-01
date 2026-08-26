package dedup

import "github.com/cespare/xxhash/v2"

func Key(channelID string, addr uint16, seq uint32) uint64 {
	digest := xxhash.New()
	digest.Write([]byte(channelID))
	digest.Write([]byte{byte(addr), byte(addr >> 8)})
	digest.Write([]byte{byte(seq), byte(seq >> 8), byte(seq >> 16), byte(seq >> 24)})
	return digest.Sum64()
}

func KeyString(channelID string, addr uint16, seq uint32) string {
	return string(rune(Key(channelID, addr, seq)))
}
