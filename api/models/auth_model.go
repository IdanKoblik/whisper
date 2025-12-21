package models

type AuthModel struct {
	ApiToken string   `bson:"_id" json:"_id"`
	Devices  []string `bson:"devices" json:"devices"`
}
