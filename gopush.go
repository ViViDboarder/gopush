package main

import (
	"flag"
	"fmt"
	"os"

	"./goson"
	"./pushbullet"
)

type Options struct {
	Token     string
	Message   string
	Device    string
	Push      bool
	List      bool
	SetActive bool
}

var options = Options{}

func loadArgs() {
	token := flag.String("token", "", "Your API Token")
	activeDevice := flag.String("d", "", "Set default device")

	flag.Parse()

	options.Token = *token
	options.Device = *activeDevice

	if options.Device != "" {
		options.SetActive = true
	}

	// Positional args
	if len(flag.Args()) == 0 {
		options.List = true
	} else if len(flag.Args()) == 1 {
		options.Message = flag.Args()[0]
		options.Push = true
	} else if len(flag.Args()) == 2 {
		options.Device = flag.Args()[0]
		options.Message = flag.Args()[1]
		options.Push = true
	}
}

func main() {
	loadArgs()

	config := goson.LoadConfig("gopush")

	var ok bool
	if options.Token != "" {
		config.Set("token", options.Token)
	} else {
		options.Token, ok = config.GetString("token")

		if !ok {
			fmt.Println("No token found")
			os.Exit(1)
		}
	}

	pb := pushbullet.New(options.Token)

	if options.Device == "" {
		activeDeviceIden, ok := config.GetString("activeDeviceIden")
		if ok {
			options.Device = activeDeviceIden
		}
	}

	pb.SetActiveDevice(options.Device)

	if options.SetActive {
		config.Set("activeDeviceIden", pb.ActiveDevice.Iden)
	}

	config.Write()

	if options.Push {
		pb.PushNote(options.Message)
	} else if options.List {
		devices := pb.GetDevices()
		PrintDevices(devices, pb.ActiveDevice)
	}
}

func PrintDevices(devices []pushbullet.Device, activeDevice pushbullet.Device) {
	fmt.Println("Devices:")
	var prefix string
	for _, device := range devices {
		if device.Iden == activeDevice.Iden {
			prefix = " *"
		} else {
			prefix = "  "
		}

		fmt.Println(prefix + device.Format())
	}
}
