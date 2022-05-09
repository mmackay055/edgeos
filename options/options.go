package options

import (
	"github.com/mmackay055/edgeos/edgedev"
	"github.com/mmackay055/edgeos/logging"
	flag "github.com/spf13/pflag"
)

var OutputFile string
var TLSCertCheckDisable bool
var Username string
var Password string
var Device edgedev.EdgeDevOpt
var Hostname string
var Command string
var Help bool
var LogFormat logging.LogFormatOpt

func init() {
	flag.StringVarP(&OutputFile, "output", "o", "", "Name of file to write backup to. Program will add extension")
	flag.BoolVarP(&TLSCertCheckDisable, "tlsdisable", "t", false, "Disable TLS Certificate Check")
	flag.StringVarP(&Username, "username", "u", "", "Set username to login to device, program will ask for it if not supplied")
	flag.StringVarP(&Password, "password", "p", "", "Set password to login to device, program will ask for it if not supplied")
	flag.VarP(&Device, "device", "d", "Set device type connecting to")
	flag.BoolVarP(&Help, "help", "h", false, "Print help menu")
	flag.VarP(&LogFormat, "logformat", "l", "Log output format selection such as systemd logs ingested by Loki")
	if !LogFormat.IsSet() {
		LogFormat.Set("cli")
	}
}
