package main

const (
	PROTO_UNKNOWN = 0
	BGP           = 1
	OSPF          = 2
)

type protocol struct {
	name       string
	ipVersion  int
	proto      int
	up         int
	imported   int64
	exported   int64
	uptime     int
	attributes map[string]interface{}
}
