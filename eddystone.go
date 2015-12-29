package main

import (
	"errors"
	"fmt"

	"github.com/paypal/gatt"
)

type frameType int

const (
	ftUID = 0x00
	ftURL = 0x10
	ftTLM = 0x20
)

type EddystoneUIDFrameField struct {
	frameType byte
	txPower   byte
	beaconID  [16]byte
	RFU       [2]byte //might not exist in some beacon
}

type EddystoneURLFrameField struct {
	frameType  byte
	txPower    byte
	urlScheme  byte
	encodedURL [18]byte
}

type EddystoneParser struct {
	BeaconName      string
	frameType       int
	FrameTypeString string
	TxPower         int8

	//UID frame only
	uidRawData EddystoneUIDFrameField
	UidString  string
	UidRFU     string

	//URL frame only
	urlRawData EddystoneURLFrameField
	UrlString  string
}

func NewEddystoneParser(adData *gatt.Advertisement) *EddystoneParser {
	ed := new(EddystoneParser)
	ed.BeaconName = adData.LocalName
	ed.frameType = int(adData.ServiceData[0].Data[0])
	ed.TxPower = int8(adData.ServiceData[0].Data[1])

	switch ed.frameType {
	case ftUID:
		fmt.Println("It is UID beancon, parse data..")
		ed.FrameTypeString = "UID"
		ed.parseUID(adData.ServiceData[0].Data)
	case ftURL:
		fmt.Println("It is URL beacon, parse data..")
		ed.FrameTypeString = "URL"
		ed.parseURL(adData.ServiceData[0].Data)
	default:
		fmt.Println("Eddystone beacon not support.")
	}

	return ed
}

func (e *EddystoneParser) parseURL(beaconData []byte) error {
	fmt.Println("beacon data size=", len(beaconData))
	if len(beaconData) != 18 && len(beaconData) != 20 {
		errString := fmt.Sprintf("Size not support uid frame:", len(beaconData), beaconData)
		return errors.New(errString)
	}

	e.urlRawData.frameType = beaconData[0]
	e.urlRawData.txPower = beaconData[1]
	return nil
}

func (e *EddystoneParser) parseUID(beaconData []byte) error {
	fmt.Println("beacon data size=", len(beaconData))
	if len(beaconData) != 18 && len(beaconData) != 20 {
		errString := fmt.Sprintf("Size not support uid frame:", len(beaconData), beaconData)
		return errors.New(errString)
	}

	//Copy data to raw data
	fmt.Printf("==> %x\n", beaconData)
	e.uidRawData.frameType = beaconData[0]
	e.uidRawData.txPower = beaconData[1]

	for i := 0; i < 16; i++ {
		e.uidRawData.beaconID[i] = beaconData[2+i]
	}

	if len(beaconData) == 20 {
		e.uidRawData.RFU[0] = beaconData[18]
		e.uidRawData.RFU[1] = beaconData[19]
	}

	//Parse UID
	//Format "%x%x%x%x-%x%x%x%x-%x%x%x%x-%x%x%x%x")
	fmt.Printf("==< %x\n", e.uidRawData.beaconID)
	for i := 0; i < 16; i++ {
		e.UidString = e.UidString + fmt.Sprintf("%x", e.uidRawData.beaconID[i])
		if (i+1)%4 == 0 {
			e.UidString = e.UidString + "-"
		}
	}
	return nil
}

func (e *EddystoneParser) PrintBeacon() {
	fmt.Println("Beacon Name:", e.BeaconName)
	fmt.Println("It is ", e.FrameTypeString, " frame eddystone")
	fmt.Println("TxPower:", e.TxPower)

	switch e.frameType {
	case ftUID:
		e.printUID()
	case ftURL:
		e.printURL()
	default:
		fmt.Println("Cannot find frame type")
	}
}

func (e *EddystoneParser) printUID() {
	fmt.Printf("UID:%x\n", e.UidString)
	if e.UidRFU != "" {
		fmt.Printf("RFU:%x\n", e.UidRFU)
	}
}

func (e *EddystoneParser) printURL() {
	fmt.Println("URL:")
}
