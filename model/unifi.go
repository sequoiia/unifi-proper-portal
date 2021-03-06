package model

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type UniFiCallbackGuest struct {
	ClientMacAddress      string
	AccessPointMacAddress string
	Timestamp             int64
	RedirectUrl           string
	Ssid                  string
}

type UniFiLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UniFiGuestAuthoriseRequest struct {
	Mac     string   `json:"mac"`
	Minutes float64  `json:"minutes"`
	Up      *float64 `json:"up,omitempty"`
	Down    *float64 `json:"down,omitempty"`
	Bytes   *float64 `json:"bytes,omitempty"`
}

type UniFiVoucherResponse struct {
	Id             string   `json:"_id"`
	AdminName      string   `json:"admin_name"`
	Code           string   `json:"code"`
	CreateTime     float64  `json:"create_time"`
	Duration       float64  `json:"duration"`
	Note           *string  `json:"note"`
	QosOverwrite   bool     `json:"qos_overwrite"`
	QosRateMaxDown *float64 `json:"qos_rate_max_down"`
	QosRateMaxUp   *float64 `json:"qos_rate_max_up"`
	QosUsageQuota  *float64 `json:"qos_usage_quota"`
	Quota          float64  `json:"quota"`
	Status         string   `json:"status"`
	StatusExpires  float64  `json:"status_expires"`
	Used           float64  `json:"used"`
}

type UniFiGuestUnauthoriseRequest struct {
	Mac string `json:"mac"`
}

func GetUniFiGuestFromCallback(r *http.Request) UniFiCallbackGuest {
	u := UniFiCallbackGuest{}
	u.ClientMacAddress = r.URL.Query().Get("id")
	u.AccessPointMacAddress = r.URL.Query().Get("ap")

	var err error
	u.Timestamp, err = strconv.ParseInt(r.URL.Query().Get("t"), 10, 64)
	if err != nil {
		u.Timestamp = 0
		log.Println(err)
	}

	u.RedirectUrl = r.URL.Query().Get("url")
	u.Ssid = r.URL.Query().Get("ssid")

	return u
}

func GetUniFiGuestCookies(r *http.Request) (UniFiCallbackGuest, error) {
	u := UniFiCallbackGuest{}
	var err error

	u.ClientMacAddress, err = GetCookieValue(r, "UPP_clientmac")
	if err != nil {
		return u, err
	}

	u.AccessPointMacAddress, err = GetCookieValue(r, "UPP_apmac")
	if err != nil {
		return u, err
	}

	u.RedirectUrl, err = GetCookieValue(r, "UPP_redirecturl")
	if err != nil {
		return u, err
	}

	u.Ssid, err = GetCookieValue(r, "UPP_ssid")
	if err != nil {
		return u, err
	}

	timestamp, err := GetCookieValue(r, "UPP_timestamp")
	if err != nil {
		return u, err
	}
	u.Timestamp, err = strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return u, err
	}

	return u, nil
}

func printAllCookies(r *http.Request) {
	cookies := r.Cookies()

	for i := 0; i < len(cookies); i++ {
		fmt.Printf("%s\n", cookies[i].String())
	}
}

func GetCookieValue(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
