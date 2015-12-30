package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/paypal/gatt"
)

type frameType int

const (
	ftUID = 0x00
	ftURL = 0x10
	ftTLM = 0x20
)

var urlSchemePrefix = []string{
	"http://www.",
	"https://www.",
	"http://",
	"https://",
}

var urlEncoding = []string{
	".com/",
	".org/",
	".edu/",
	".net/",
	".info/",
	".biz/",
	".gov/",
	".com",
	".org",
	".edu",
	".net",
	".info",
	".biz",
	".gov",
}

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
	encodedURL []byte
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

func decodeURL(prefix byte, encodedURL []byte) (string, error) {
	if int(prefix) >= len(urlSchemePrefix) {
		return "", errors.New("invaild prefix")
	}

	s := urlSchemePrefix[prefix]

	for _, b := range encodedURL {
		switch {
		case 0x00 <= b && b <= 0x13:
			s += urlEncoding[b]
		case 0x0e <= b && b <= 0x20:
			fallthrough
		case 0x7f <= b && b <= 0xff:
			return "", errors.New("invalid byte")
		default:
			s += string(b)
		}
	}

	return s, nil
}

func (e *EddystoneParser) parseURL(beaconData []byte) error {
	//Copy raw data first
	e.urlRawData.frameType = beaconData[0]
	e.urlRawData.txPower = beaconData[1]
	e.urlRawData.urlScheme = beaconData[2]
	for i := 0; i+3 < len(beaconData); i++ {
		e.urlRawData.encodedURL = append(e.urlRawData.encodedURL, beaconData[3+i])
	}
	//Decode URL data
	url, _ := decodeURL(e.urlRawData.urlScheme, e.urlRawData.encodedURL)
	e.UrlString = url
	return nil
}

func (e *EddystoneParser) parseUID(beaconData []byte) error {
	if len(beaconData) != 18 && len(beaconData) != 20 {
		errString := fmt.Sprintf("Size not support uid frame:", len(beaconData), beaconData)
		return errors.New(errString)
	}

	//Copy data to raw data
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
	e.UidString = fmt.Sprintf("%x-%x-%x-%x", e.uidRawData.beaconID[:4], e.uidRawData.beaconID[4:8], e.uidRawData.beaconID[8:12], e.uidRawData.beaconID[12:16])
	e.UidString = strings.ToUpper(e.UidString)
	return nil
}

//Print the beacon detail information
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
	fmt.Printf("UID: %s\n", e.UidString)
	if e.UidRFU != "" {
		fmt.Printf("RFU:%x\n", e.UidRFU)
	}
}

func (e *EddystoneParser) printURL() {
	fmt.Println("URL:", e.UrlString)
}
