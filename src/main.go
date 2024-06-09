package main

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"time"

	"tinygo.org/x/bluetooth"
)

func main() {
	adapter := bluetooth.DefaultAdapter

	secretEnv := os.Getenv("SECRET")
	if secretEnv == "" {
		slog.Error("SECRET environment variable not set")
		return
	}
	secret := []byte(secretEnv)

	timeoutEnv := os.Getenv("TIMEOUT")
	if timeoutEnv == "" {
		slog.Error("TIMEOUT environment variable not set")
		return
	}
	timeout, err := time.ParseDuration(timeoutEnv)
	if err != nil {
		slog.Error("TIMEOUT environment variable parse error", err)
		return
	}

	slog.Info("enable bluetooth")
	err = adapter.Enable()
	if err != nil {
		slog.Error("enable bluetooth error", err)
		return
	}

	slog.Info("start scanning")
	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		b, ok := getBeacon(device)
		if !ok {
			return
		}
		ok = checkSignature(b, secret)
		if !ok {
			return
		}
		ok = isNewIBeacon(b, timeout)
		if !ok {
			return
		}
		fmt.Printf("\n%+v\n", b)
	})
	if err != nil {
		slog.Error("scan error", err)
	}
}

func isIBeacon(manData []bluetooth.ManufacturerDataElement) ([]byte, bool) {
	if len(manData) == 0 {
		return nil, false
	}
	if manData[0].CompanyID != 0x004C {
		return nil, false
	}
	ok := manData[0].Data[0] == 0x02 && manData[0].Data[1] == 0x15
	if !ok {
		return nil, false
	}
	return manData[0].Data, true
}

func getBeacon(device bluetooth.ScanResult) (res *IBeacon, ok bool) {
	manData := device.ManufacturerData()
	byteArray, ok := isIBeacon(manData)
	if !ok {
		return
	}

	res = &IBeacon{
		Time:             time.Now(),
		ScanResult:       device,
		ManufacturerData: manData,
		Address:          device.Address,
		RSSI:             device.RSSI,
		LocalName:        device.LocalName(),
		SubType:          byteArray[0],
		SubTypeLength:    byteArray[1],
		UMMID:            byteArray[2:22],
		UUID:             byteArray[2:18],
		Major:            binary.BigEndian.Uint16(byteArray[18:20]),
		Minor:            binary.BigEndian.Uint16(byteArray[20:22]),
		MeasuredPower:    int8(byteArray[22]),
	}

	return
}

type IBeacon struct {
	Time             time.Time
	ScanResult       bluetooth.ScanResult
	ManufacturerData []bluetooth.ManufacturerDataElement
	Address          bluetooth.Address
	RSSI             int16
	LocalName        string
	SubType          uint8
	SubTypeLength    uint8
	UMMID            []byte
	UUID             []byte
	Major            uint16
	Minor            uint16
	MeasuredPower    int8
}
