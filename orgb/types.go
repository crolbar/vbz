package orgb

import "net"

const (
	NET_PACKET_ID_REQUEST_CONTROLLER_COUNT      = 0
	NET_PACKET_ID_REQUEST_CONTROLLER_DATA       = 1
	NET_PACKET_ID_REQUEST_PROTOCOL_VERSION      = 40
	NET_PACKET_ID_SET_CLIENT_NAME               = 50
	NET_PACKET_ID_DEVICE_LIST_UPDATED           = 100
	NET_PACKET_ID_REQUEST_PROFILE_LIST          = 150
	NET_PACKET_ID_REQUEST_SAVE_PROFILE          = 151
	NET_PACKET_ID_REQUEST_LOAD_PROFILE          = 152
	NET_PACKET_ID_REQUEST_DELETE_PROFILE        = 153
	NET_PACKET_ID_REQUEST_PLUGIN_LIST           = 200
	NET_PACKET_ID_PLUGIN_SPECIFIC               = 201
	NET_PACKET_ID_RGBCONTROLLER_RESIZEZONE      = 1000
	NET_PACKET_ID_RGBCONTROLLER_CLEARSEGMENTS   = 1001
	NET_PACKET_ID_RGBCONTROLLER_ADDSEGMENT      = 1002
	NET_PACKET_ID_RGBCONTROLLER_UPDATELEDS      = 1050
	NET_PACKET_ID_RGBCONTROLLER_UPDATEZONELEDS  = 1051
	NET_PACKET_ID_RGBCONTROLLER_UPDATESINGLELED = 1052
	NET_PACKET_ID_RGBCONTROLLER_SETCUSTOMMODE   = 1100
	NET_PACKET_ID_RGBCONTROLLER_UPDATEMODE      = 1101
	NET_PACKET_ID_RGBCONTROLLER_SAVEMODE        = 1102
)

type ORGBConn struct {
	conn net.Conn
}

type ORGBHeader struct {
	devIdx uint32
	id     uint32
	size   uint32
}

type Controller struct {
	Type        int
	Name        string
	Description string
	Version     string
	Serial      string
	Location    string
	ActiveMode  int
	Colors      []RGBColor

	// uneeded
	// Modes       []Mode
	// Zones       []Zone
	// Leds        []LED
}

type RGBColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
}
