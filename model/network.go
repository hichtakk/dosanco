package model

import (
	"strings"
	"strconv"
	"net"
)

type IPv4Network struct {
	Model
	CIDR 			string			`json:"cidr" validate:"required" gorm:"unique;not null"`
	Description		string			`json:"description"`
	//Supernetwork	*IPv4Network
	SupernetworkID	uint			`json:"supernet_id" validate:"required"`
	//Subnetwork		[]*IPv4Network	`gorm:"many2many:ipv4_subnetwork;association_jointable_foreignkey:subnet_id"`
}

func (n IPv4Network) GetNetwork() *net.IPNet {
	_, ipv4Net, _ := net.ParseCIDR(n.CIDR)
	return ipv4Net
}

func (n IPv4Network) GetPrefixLength() int {
	slice := strings.Split(n.CIDR, "/")
	len, _ := strconv.Atoi(slice[1])
	return len
}