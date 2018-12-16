package mrbluetooth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

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

// FlipRequest is schema of payload for flip action
type FlipRequest struct {
	Seq []struct {
		X    uint8 `json:"x"`
		Y    uint8 `json:"y"`
		Pole int   `json:"pole"`
	} `json:"seq"`
	Interval uint32 `json:"interval"` // in milli seconds
}

func (f *FlipRequest) Parse(data []byte) error {
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
			rsp.Write(payload.Compose())
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
		gatt.CentralConnected(func(c gatt.Central) { log.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { log.Println("Disconnect: ", c.ID()) }),
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
