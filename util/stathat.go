package util

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var statHatKey string

func StatHatSetKey(key string) {
	statHatKey = key
	fmt.Printf("STATHATSETKEY: %s\n", key)
}

func StatHatCount(stat string, count float32) error {
	if statHatKey == "" {
		fmt.Println("STATHATCOUNT: No StatHat key has been supplied via util.StatHatSetKey")
		return errors.New("No StatHat key has been supplied via util.StatHatSetKey")
	}
	fmt.Printf("STATHATCOUNT: %s = %f\n", stat, count)
	r, err := http.PostForm("https://api.stathat.com/ez", url.Values{
		"ezkey": {statHatKey},
		"stat":  {stat},
		"count": {fmt.Sprintf("%f", count)},
	})
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func StatHatValue(stat string, value float32) error {
	if statHatKey == "" {
		fmt.Println("STATHATVALUE: No StatHat key has been supplied via util.StatHatSetKey")
		return errors.New("No StatHat key has been supplied via util.StatHatSetKey")
	}
	fmt.Printf("STATHATVALUE: %s = %f\n", stat, value)
	r, err := http.PostForm("https://api.stathat.com/ez", url.Values{
		"ezkey": {statHatKey},
		"stat":  {stat},
		"value": {fmt.Sprintf("%f", value)},
	})
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}
