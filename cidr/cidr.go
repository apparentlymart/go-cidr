// Package cidr is a collection of assorted utilities for computing
// network and host addresses within network ranges.
//
// It expects a CIDR-type address structure where addresses are divided into
// some number of prefix bits representing the network and then the remaining
// suffix bits represent the host.
//
// For example, it can help to calculate addresses for sub-networks of a
// parent network, or to calculate host addresses within a particular prefix.
//
// At present this package is prioritizing simplicity of implementation and
// de-prioritizing speed and memory usage. Thus caution is advised before
// using this package in performance-critical applications or hot codepaths.
// Patches to improve the speed and memory usage may be accepted as long as
// they do not result in a significant increase in code complexity.
package cidr

import (
	"fmt"
	"net"
)

// Subnet takes a parent CIDR range and creates a subnet within it
// with the given number of additional prefix bits and the given
// network number.
//
// For example, 10.3.0.0/16, extended by 8 bits, with a network number
// of 5, becomes 10.3.5.0/24 .
func Subnet(base *net.IPNet, newBits int, num int) (*net.IPNet, error) {
	ip := base.IP
	mask := base.Mask

	parentLen, addrLen := mask.Size()
	newPrefixLen := parentLen + newBits

	if newPrefixLen > addrLen {
		return nil, fmt.Errorf("insufficient address space to extend prefix of %d by %d", parentLen, newBits)
	}

	maxNetNum := uint64(1<<uint64(newBits)) - 1
	if uint64(num) > maxNetNum {
		return nil, fmt.Errorf("prefix extension of %d does not accommodate a subnet numbered %d", newBits, num)
	}

	return &net.IPNet{
		IP:   insertNumIntoIP(ip, num, newPrefixLen),
		Mask: net.CIDRMask(newPrefixLen, addrLen),
	}, nil
}

// Host takes a parent CIDR range and turns it into a host IP address with
// the given host number.
//
// For example, 10.3.0.0/16 with a host number of 2 gives 10.3.0.2.
func Host(base *net.IPNet, num int) (net.IP, error) {
	ip := base.IP
	mask := base.Mask

	parentLen, addrLen := mask.Size()
	hostLen := addrLen - parentLen

	maxHostNum := uint64(1<<uint64(hostLen)) - 1
	if uint64(num) > maxHostNum {
		return nil, fmt.Errorf("prefix of %d does not accommodate a host numbered %d", parentLen, num)
	}

	return insertNumIntoIP(ip, num, 32), nil
}
