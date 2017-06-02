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

const (
	action_list_devs = iota
	action_send_note = iota
	action_send_link = iota
)

type Options struct {
	Token   string
	Message string
	Device  string
	Action  int
}

func readOptions() (options Options, err error) {
	tokenPtr := flag.String("token", "", "Your API Token (will be persisted to your home directory)")
	devicePtr := flag.String("d", "", "Set default device (defaults to all devices)")
	listDevices := flag.Bool("l", false, "List devices")
	pushType := flag.String("t", "note", "Push type (note or link)")

	flag.Parse()

	// Token is needed for any connection to PushBullet
	options = Options{}
	options.Token = *tokenPtr
	options.Device = *devicePtr

	if *listDevices {
		options.Action = action_list_devs
		return
	}

	switch *pushType {
	case "link":
		option.Action = action_send_link
		break
	default:
		// TODO: Decide if this should warn and print usage
	case "note":
		options.Action = action_send_note
		break
	}

	// Read message from stdin or args
	if flag.NArg() == 0 {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return options, err
		}
		options.Message = string(data)
	} else {
		options.Message = strings.Join(flag.Args(), " ")
	}

	if option.Message == "" {
		options.Action = action_list_devs
	}
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
	options, err := readOptions()
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

	switch option.Action {
	case action_send_note:
		pb.PushNote(options.Message)
		break
	default:
	case action_list_devs:
		devices := pb.GetDevices()
		printDevices(devices, pb.ActiveDevice)
		break
	}
}
