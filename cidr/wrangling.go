package cidr

import (
	"fmt"
	"math/big"
	"net"
)

func ipToInt(ip net.IP) (*big.Int, int) {
	if len(ip) == net.IPv4len {
		return big.NewInt(int64(int32(ip[0])<<24 | int32(ip[1])<<16 | int32(ip[2])<<8 | int32(ip[3]))), 32
	} else if len(ip) == net.IPv6len {
		bottom := big.NewInt(int64(
			uint64(ip[0])<<56 | uint64(ip[1])<<48 | uint64(ip[2])<<40 |
			uint64(ip[3])<<32 | uint64(ip[4])<<24 | uint64(ip[5])<<16 |
			uint64(ip[6])<<8 | uint64(ip[7]),
		))
		top := big.NewInt(int64(
			uint64(ip[8])<<56 | uint64(ip[9])<<48 | uint64(ip[10])<<40 |
			uint64(ip[11])<<32 | uint64(ip[12])<<24 | uint64(ip[13])<<16 |
			uint64(ip[14])<<8 | uint64(ip[15]),
		))
		bottom.Lsh(bottom, 64)
		return bottom.Or(bottom, top), 128
	} else {
		panic(fmt.Errorf("Unsupported address length %d", len(ip)))
	}
}
