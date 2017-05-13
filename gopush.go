package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	switch flag.NArg() {
	case 0:
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			// TODO: Make errors great again!
			os.Exit(1)
		}
		message := string(data)
		if message == "" {
			options.List = true
		} else {
			fmt.Println(message)
			options.Message = message
			options.Push = true
		}
		break
	case 1:
		options.Message = flag.Args()[0]
		options.Push = true
		break
	case 2:
		options.Device = flag.Args()[0]
		options.Message = flag.Args()[1]
		options.Push = true
		break
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
