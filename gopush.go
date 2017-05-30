package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ViViDboarder/gopush/goson"
	"github.com/ViViDboarder/gopush/pushbullet"
)

type Options struct {
	Token   string
	Message string
	Device  string
	Push    bool
	List    bool
}

func loadArgs() (options Options, err error) {
	options = Options{}
	token := flag.String("token", "", "Your API Token")
	activeDevice := flag.String("d", "", "Set default device")
	listDevices := flag.Bool("l", false, "List devices")

	flag.Parse()

	options.Token = *token
	options.Device = *activeDevice
	options.List = *listDevices

	if options.List {
		return
	}

	// Positional args
	switch flag.NArg() {
	case 0:
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return options, err
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
	return
}

func printDevices(devices []pushbullet.Device, activeDevice pushbullet.Device) {
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

func getToken(options Options, config goson.Config) (token string, err error) {
	if options.Token != "" {
		token = options.Token
		config.Set("token", options.Token)
	} else {
		var ok bool
		token, ok = config.GetString("token")
		if !ok {
			err = errors.New("No token found")
		}
	}
	return
}

func getDevice(options Options, config goson.Config) (device string) {
	if options.Device == "" {
		activeDeviceIden, ok := config.GetString("activeDeviceIden")
		if ok {
			device = activeDeviceIden
		}
	}
	return
}

func persistDevice(activeDeviceIden string, config goson.Config) {
	storedDeviceIden, ok := config.GetString("activeDeviceIden")
	if activeDeviceIden != "" && ok && activeDeviceIden != storedDeviceIden {
		config.Set("activeDeviceIden", activeDeviceIden)
	}
}

func main() {
	options, err := loadArgs()
	if err != nil {
		log.Fatal(err)
	}
	config := goson.LoadConfig("gopush")
	token, err := getToken(options, config)
	if err != nil {
		fmt.Println("Token required and not provided by config or command line")
		log.Fatal(err)
	}
	device := getDevice(options, config)

	pb := pushbullet.New(token)
	pb.SetActiveDevice(device)

	persistDevice(pb.ActiveDevice.Iden, config)
	config.Write()

	if options.Push {
		pb.PushNote(options.Message)
	} else if options.List {
		devices := pb.GetDevices()
		printDevices(devices, pb.ActiveDevice)
	}
}
