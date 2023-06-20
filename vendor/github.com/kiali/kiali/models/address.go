package models

import core_v1 "k8s.io/api/core/v1"

type Addresses []Address
type Address struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

func (addresses *Addresses) Parse(as []core_v1.EndpointAddress) {
	for _, address := range as {
		castedAddress := Address{}
		castedAddress.Parse(address)
		*addresses = append(*addresses, castedAddress)
	}
}

func (address *Address) Parse(a core_v1.EndpointAddress) {
	address.IP = a.IP

	if a.TargetRef != nil {
		address.Kind = a.TargetRef.Kind
		address.Name = a.TargetRef.Name
	}
}
