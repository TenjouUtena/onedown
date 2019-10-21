package session

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestSessionMessage_Unmarshal(t *testing.T) {
	messagePayload := "{\"name\":\"WriteSquare\",\"session\":\"00000000-0000-0000-0000-000000000000\",\"payload\":" +
		"\"{\\\"solver\\\":\\\"00000000-0000-0000-0000-000000000000\\\",\\\"row\\\":0,\\\"col\\\":0,\\\"answer\\\":\\\"Q\\\"}\"" +
		"}"
	messagePayloadBytes := []byte(messagePayload)
	expected := MessageForSession{
		SessionId: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Message: WriteSquare{
			Solver: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			Row:    0,
			Col:    0,
			Answer: "Q",
		},
	}
	message, err := unmarshallSocketMessage(messagePayloadBytes)
	if err != nil {
		t.Error(err)
	} else if message != expected {
		fmt.Printf("Actual did not match expected:\n")
		fmt.Printf("ACTUAL:   %s\n", message)
		fmt.Printf("EXPECTED: %s\n", expected)
		t.Fail()
	}
}
