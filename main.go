package main

import (
	"bufio"
	"github.com/mmackay055/edgeos/options"
	"github.com/mmackay055/edgeos/edgedev"
	"fmt"
	flag "github.com/spf13/pflag"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

const (
	ERR_HOSTNAME int = 2
	ERR_LOGIN    int = 3
	ERR_DEVICE   int = 4
	ERR_COMMAND  int = 5
        ERR_BACKUP   int = 6
        ERR_SAVE     int = 7
)

func main() {
	flag.Parse()

	if options.Help {
		printHelp()
		os.Exit(0)
	}

	// request username if not supplied
	if len(options.Username) == 0 {
		options.Username = getUsername()
	}

	// request password if not supplied
	if len(options.Password) == 0 {
		options.Password = getPassword()
	}

	// request device type
	if !options.Device.IsSet() {
		err := options.Device.Set(getDevice())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(ERR_DEVICE)
		}
	}

	args := flag.Args()

	// get hostname
	if len(args) < 1 {
		errorFound("error require device hostname/IP", ERR_HOSTNAME)
	}
	options.Hostname = args[0]

	// get command
	if len(args) < 2 {
		errorFound("error require command", ERR_COMMAND)
	}
	options.Command = args[1]

	//select device
	device := options.Device.GetDevice()
        edgedev.DisableTLSCertCheck(device, options.TLSCertCheckDisable)

	// Login to device
	if err := edgedev.Login(device, 
                options.Hostname, 
                options.Username, 
                options.Password); err != nil {

		errorFound(fmt.Sprintf("error login: %s", err),
                        ERR_LOGIN)
	}

        // execute command
	switch options.Command {
	case "backup":
		if err := edgedev.BackUp(device, 
                                        options.OutputFile); err != nil {

			errorFound(fmt.Sprintf("error backup: %s", err),
                                                ERR_BACKUP)
		}
        case "save":
                if err := edgedev.Save(device); err != nil {
			errorFound(fmt.Sprintf("error save: %s", err),
                                  ERR_SAVE)
		}
	default:
		errorFound("error invalid command",
                        ERR_COMMAND)
	}

}

func printHelp() {
	fmt.Fprintf(os.Stderr, "Command line tool for performing operations on Ubiquiti Edge Devices.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "\tedgeoscli [options] hostname command\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "\tsave - save running configuration on device.\n")
	fmt.Fprintf(os.Stderr, "\tbackup - backup saved configuration on device.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func getUsername() string {
	stdin := bufio.NewReader(os.Stdin)

	fmt.Print("username:")
	username, err := stdin.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(username)
}

func getPassword() string {
	fmt.Print("password:")
	passByte, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	fmt.Println()

	return strings.TrimSpace(string(passByte))
}

func getDevice() string {
	stdin := bufio.NewReader(os.Stdin)

	fmt.Printf("device(%v):", edgedev.AvailableEdgeDevs())

	device, err := stdin.ReadString('\n')
	if err != nil {
		return ""
	}

	return strings.TrimSpace(device)
}

func errorFound(msg string, exitStatus int) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(exitStatus)
}
