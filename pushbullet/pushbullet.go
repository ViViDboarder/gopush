package pushbullet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	host = "https://api.pushbullet.com/v2"
)

func New(apiKey string) PushBullet {
	return PushBullet{apiKey: apiKey}
}

type PushBullet struct {
	apiKey       string
	client       *http.Client
	ActiveDevice Device
	Devices      []Device
}

// Loads all devices for account into instance
func (pb *PushBullet) LoadDevices() []Device {
	dl := new(DeviceList)
	err := pb.pbRequest("GET", host+"/devices", nil, dl)

	if err != nil {
		fmt.Println("error on pbRequest")
		fmt.Println(err)
	}

	pb.Devices = dl.Devices

	return dl.Devices
}

// Lazy loaded version of LoadDevices
func (pb *PushBullet) GetDevices() []Device {
	if len(pb.Devices) == 0 {
		pb.LoadDevices()
	}

	return pb.Devices
}

func (pb PushBullet) PushNote(message string) {
	body := map[string]interface{}{
		"type": "note",
		"body": message,
	}

	pb.Push(body)
}

func (pb PushBullet) PushLink(title string, url string) {
	pb.PushLinkWithBody(title, url, "")
}

func (pb PushBullet) PushLinkWithBody(title string, url string, messBody string) {
	body := map[string]interface{}{
		"type":  "link",
		"title": title,
		"url":   url,
		"body":  messBody,
	}

	pb.Push(body)
}

// Pushes a text message to active device
// TODO: handle response
func (pb PushBullet) Push(body map[string]interface{}) {
	body["device_iden"] = pb.ActiveDevice.Iden

	result := new(interface{})
	pb.pbRequest("POST", host+"/pushes", body, result)
}

// Sets active device by Nickname or Iden
func (pb *PushBullet) SetActiveDevice(key string) {
	if len(pb.Devices) == 0 {
		pb.LoadDevices()
	}

	for _, d := range pb.Devices {
		if d.Match(key) && d.CanPush() {
			pb.ActiveDevice = d
		}
	}
}

// Method for communicating with PushBullet
func (pb *PushBullet) pbRequest(method string, url string, body map[string]interface{}, result interface{}) (err error) {
	if pb.client == nil {
		pb.client = &http.Client{}
	}

	var r *http.Request
	if body == nil {
		r, err = http.NewRequest(method, url, nil)
	} else {
		jsonBody, jsonError := json.Marshal(body)
		if jsonError != nil {
			fmt.Println("Error marshalling body: " + url)
			return jsonError
		}
		r, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	}

	if err != nil {
		fmt.Println("Error on request create: " + url)
		return
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Application-Type", "application/json")
	r.Header.Add("User-Agent", "gopush")

	r.SetBasicAuth(pb.apiKey, "")

	resp, err := pb.client.Do(r)
	if err != nil {
		fmt.Println("Error on request do: ")
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error on response read: " + url)
		return
	}

	err = json.Unmarshal(responseBody, result)
	if err != nil {
		fmt.Println("Error on response parse: " + url)
		fmt.Println("json: " + string(responseBody))
		return
	}

	return err
}

// Response types

type DeviceList struct {
	Aliases []interface{}
	Clients []interface{}
	Devices []Device
	Grants  []interface{}
	Pushes  []interface{}
}

type DeviceFingerprint struct {
	mac_address string
	android_id  string
}

type Device struct {
	Iden         string
	Created      float32
	Nickname     string
	Modified     float32
	Push_token   string
	Fingerprint  string
	Active       bool
	Model        string
	App_version  int32
	Type         string
	Kind         string
	Pushable     bool
	Manufacturer string
}

// Returns Android Fingerprint
func (d Device) AndroidFingerprint() (DeviceFingerprint, error) {
	fp := new(DeviceFingerprint)
	err := json.Unmarshal([]byte(d.Fingerprint), fp)
	return *fp, err
}

// Returns if the device is available for pushing
func (d Device) CanPush() bool {
	return d.Active && d.Pushable
}

// Outputs the device in a friendly string format
func (d Device) Format() string {
	return fmt.Sprintf("%s (%s)", d.Nickname, d.Kind)
}

// Returns true if key matches Nickname or Iden
func (d Device) Match(key string) bool {
	if d.Nickname == key || d.Iden == key {
		return true
	}

	return false
}
