package models

type MQTTUser struct {
	CliendID string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Protocol string   `json:"mqttprotocol"`
	Host     string   `json:"mqtthost"`
	Port     int32    `json:"mqttport"`
	User     MQTTUser `json:"mqttuser"`
}
