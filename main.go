package main

import (
	"flag"
	"log"
	"time"

	"github.com/kotomiya/co2mini/mackerel"
	"github.com/zserge/hid"
)

const (
	vendorID = "04d9:a052:0100:00"
	co2op    = 0x50
	tempop   = 0x42
)

var (
	key = []byte{0x86, 0x41, 0xc9, 0xa8, 0x7f, 0x41, 0x3c, 0xac}
)

type metric struct {
	Type  string
	Value int
}

func monitor(device hid.Device, c *mackerel.Client) {
	if err := device.Open(); err != nil {
		log.Println("Open error: ", err)
		return
	}
	defer device.Close()

	if err := device.SetReport(0, key); err != nil {
		log.Fatal(err)
	}

	for {
		buf, err := device.Read(-1, 1*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		dec := decrypt(buf, key)
		if len(dec) == 0 {
			continue
		}
		val := int(dec[1])<<8 | int(dec[2])
		if dec[0] == co2op {
			c.Post(mackerel.Metric{{Name: "co2", Time: time.Now().Unix(), Value: val}})
			log.Printf("co2:%d ppm", val)
		}
		if dec[0] == tempop {
			temp := int(float64(val)/16.0 - 273.15)
			c.Post(mackerel.Metric{{Name: "temp", Time: time.Now().Unix(), Value: temp}})
			log.Printf("temp:%d ", temp)
		}
	}
}

func main() {
	flag.Parse()
	c := mackerel.NewClient()
	hid.UsbWalk(func(device hid.Device) {
		monitor(device, &c)
	})
}
