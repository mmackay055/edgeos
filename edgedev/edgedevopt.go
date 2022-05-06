package edgedev

import (
	"errors"
	"fmt"
)

var devices []string

const (
        edgeRouter = `edgerouter`  
	edgeSwitch = `edgeswitch`
)

// EdgeDevOpt option definition
type EdgeDevOpt struct {
	name string
	dev  EdgeDev
}

func (d *EdgeDevOpt) String() string {
	return d.name
}

func (d *EdgeDevOpt) Type() string {
	return "Device"
}

func (d *EdgeDevOpt) Set(val string) error {
	if len(val) < 1 {
		return nil
	}
	return d.setDevice(val)
}

func AvailableEdgeDevs() []string {
	available := [2]string{edgeRouter, edgeSwitch}
        return available[:]
}

func (d *EdgeDevOpt) setDevice(devname string) error {
	var dev EdgeDev

	switch devname {
	case edgeRouter:
		dev = &EdgeRouter{}
		d.name = edgeRouter
	case edgeSwitch:
		dev = &EdgeSwitch{}
		d.name = edgeSwitch 
	default:
		return errors.New(fmt.Sprintf("device %s is not suppored",
			devname))
	}
	d.dev = dev
	return nil
}

func (d EdgeDevOpt) GetDevice() EdgeDev {
	return d.dev
}

func (d *EdgeDevOpt) IsSet() bool {
	return d.dev != nil
}
