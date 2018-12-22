package mrbluetooth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/service"
	"github.com/paypal/gatt/linux/cmd"
)

const (
	device_name = "MagicReversi"
)

type cmdReadBDAddr struct{}

func (c cmdReadBDAddr) Marshal(b []byte) {}
func (c cmdReadBDAddr) Opcode() int      { return 0x1009 }
func (c cmdReadBDAddr) Len() int         { return 0 }

func bdaddr(d gatt.Device) {
	rsp := bytes.NewBuffer(nil)
	if err := d.Option(gatt.LnxSendHCIRawCommand(&cmdReadBDAddr{}, rsp)); err != nil {
		fmt.Printf("Failed to send HCI raw command, err: %s", err)
	}
	b := rsp.Bytes()
	if b[0] != 0 {
		fmt.Printf("Failed to get bdaddr with HCI Raw command, status: %d", b[0])
	}
	log.Printf("BD Addr: %02X:%02X:%02X:%02X:%02X:%02X", b[6], b[5], b[4], b[3], b[2], b[1])
}

type Payload interface {
	// Compose returns JSON string as byte data
	Compose() []byte
}

// FlipEvent represents a flip event
type FlipEvent struct {
	X    uint8
	Y    uint8
	Pole int
}

// FlipRequest is schema of payload for flip action
type FlipRequest struct {
	Seq      []FlipEvent
	Interval time.Duration
}

func (f *FlipRequest) Parse(data []byte) error {
	/*
		Flip Protocol
		| interval (8bit) | padding(1bit) pole(1bit) x(3bit)  y(3bit) | … |
		First 1 byte means interval time (10msec) between flips
	*/

	f.Interval = time.Duration(data[0]) * 10 * time.Millisecond
	f.Seq = []FlipEvent{}

	for _, b := range data[1:] {

		f.Seq = append(f.Seq, FlipEvent{
			Pole: int(0x40&b) >> 6,
			X:    uint8(0x38&b) >> 3,
			Y:    uint8(0x07&b) >> 0,
		})
	}

	err := json.Unmarshal(data, f)
	if err != nil {
		return err
	}

	return nil
}

func NewMagicReversiService(
	FlipHandle func(FlipRequest),
	GetBoardHandle func() Payload,
	notify_chan chan Payload) *gatt.Service {
	s := gatt.NewService(gatt.MustParseUUID("BF3C7852-3D5D-4C4B-88C7-D1F01A448854"))

	s.AddCharacteristic(gatt.MustParseUUID("B1F94B37-53A8-453E-9C21-43A9D9F2A7E3")).HandleNotifyFunc(
		func(r gatt.Request, n gatt.Notifier) {
			for payload := range notify_chan {
				msg := string(payload.Compose())
				fmt.Fprintf(n, msg)
			}
		})

	s.AddCharacteristic(gatt.MustParseUUID("1A0BF805-29DA-4AC0-BE94-FE5EDDA879D8")).HandleWriteFunc(
		func(r gatt.Request, data []byte) (status byte) {
			payload := FlipRequest{}
			err := payload.Parse(data)
			if err != nil {
				log.Println("FlipRequest: Parse Error")
			}

			FlipHandle(payload)
			return gatt.StatusSuccess
		})

	s.AddCharacteristic(gatt.MustParseUUID("183600AD-A496-46A1-95D2-79B4D1673DB1")).HandleReadFunc(
		func(rsp gatt.ResponseWriter, req *gatt.ReadRequest) {

			payload := GetBoardHandle()
			c := payload.Compose()
			_, err := rsp.Write(c)

			if err != nil {
				log.Println("error: ", err)
			}
			return
		})
	return s
}

// MrBluetooth is an object for bluetooth connection
type MrBluetooth struct{}

func NewMrBluetooth() MrBluetooth {
	return MrBluetooth{}
}

func (*MrBluetooth) Launch(
	FlipHandle func(FlipRequest),
	GetBoardHandle func() Payload,
	notify_chan chan Payload) {

	d, err := gatt.NewDevice(
		gatt.LnxMaxConnections(1),
		gatt.LnxDeviceID(-1, true),
		gatt.LnxSetAdvertisingParameters(&cmd.LESetAdvertisingParameters{
			AdvertisingIntervalMin: 0x00f4,
			AdvertisingIntervalMax: 0x00f4,
			AdvertisingChannelMap:  0x07,
		}),
	)

	if err != nil {
		log.Printf("Failed to open device, err: %s", err)
		return
	}

	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) {
			log.Printf("info: Connect: %s, MTU: %d", c.ID(), c.MTU())
		}),
		gatt.CentralDisconnected(func(c gatt.Central) {
			log.Println("Disconnect: ", c.ID())
		}),
	)

	onStateChanged := func(d gatt.Device, s gatt.State) {
		fmt.Printf("State: %s\n", s)
		switch s {
		case gatt.StatePoweredOn:
			bdaddr(d)

			d.AddService(service.NewGapService(device_name))
			d.AddService(service.NewGattService())

			mr_service := NewMagicReversiService(FlipHandle, GetBoardHandle, notify_chan)
			d.AddService(mr_service)

			// Advertise
			d.AdvertiseNameAndServices(device_name, []gatt.UUID{mr_service.UUID()})

		default:
		}
	}

	d.Init(onStateChanged)
	select {}
}
