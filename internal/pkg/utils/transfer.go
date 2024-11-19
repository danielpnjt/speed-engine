package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func GeneratePaymentRef(AccountNumber string) (string, error) {
	randomBytes := make([]byte, 5)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := strings.ToUpper(hex.EncodeToString(randomBytes))
	phoneLast5 := AccountNumber[len(AccountNumber)-5:]
	secondsInDay := time.Now().Hour()*3600 + time.Now().Minute()*60 + time.Now().Second()
	timeComponent := fmt.Sprintf("%05d", secondsInDay)
	paymentCode := fmt.Sprintf("TF-%s%s%s", randomString[:3], phoneLast5[:2], timeComponent[:5])

	return paymentCode, nil
}
