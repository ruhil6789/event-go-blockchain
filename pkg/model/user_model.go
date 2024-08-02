package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`       // MongoDB's default ID field
	WalletAddress   string             `bson:"walletAddress"`       // Wallet address (string)
	ReferralLink    string             `bson:"referralLink"`        // Referral link (string)
	CountOfReferee  int                `bson:"countOfReferree"`     // Count of referees (number)
	ReferrerAddress string             `bson:"referrerAddress"`     // Referrer address (string)
	TotalStaked     string             `bson:"totalStaked"`         // Total staked (string)
	ClaimableAmount float64            `bson:"claimableAmount"`     // Claimable amount (number)
	PackageAmount   float64            `bson:"packageAmount"`       // Package amount (number)
	Added           bool               `bson:"added"`               // Added flag (boolean)
	Timestamp       *int64             `bson:"timestamp,omitempty"` // Timestamp (number)
	IsLockUser      bool               `bson:"isLockUser"`          // Is the user locked? (boolean)
	Level           Level              `bson:"level"`               // Nested levels (object)
}

// Level structure to handle the nested array fields
type Level struct {
	Lvl0  []string `bson:"lvl0"`
	Lvl1  []string `bson:"lvl1"`
	Lvl2  []string `bson:"lvl2"`
	Lvl3  []string `bson:"lvl3"`
	Lvl4  []string `bson:"lvl4"`
	Lvl5  []string `bson:"lvl5"`
	Lvl6  []string `bson:"lvl6"`
	Lvl7  []string `bson:"lvl7"`
	Lvl8  []string `bson:"lvl8"`
	Lvl9  []string `bson:"lvl9"`
	Lvl10 []string `bson:"lvl10"`
	Lvl11 []string `bson:"lvl11"`
	Lvl12 []string `bson:"lvl12"`
	Lvl13 []string `bson:"lvl13"`
	Lvl14 []string `bson:"lvl14"`
	Lvl15 []string `bson:"lvl15"`
	Lvl16 []string `bson:"lvl16"`
}

// func NewUser() *User {
// 	return &User{
// 		Level: Level{
// 			Lvl0:  []string{},
// 			Lvl1:  []string{},
// 			Lvl2:  []string{},
// 			Lvl3:  []string{},
// 			Lvl4:  []string{},
// 			Lvl5:  []string{},
// 			Lvl6:  []string{},
// 			Lvl7:  []string{},
// 			Lvl8:  []string{},
// 			Lvl9:  []string{},
// 			Lvl10: []string{},
// 			Lvl11: []string{},
// 			Lvl12: []string{},
// 			Lvl13: []string{},
// 			Lvl14: []string{},
// 			Lvl15: []string{},
// 			Lvl16: []string{},
// 		},
// 	}
// }
