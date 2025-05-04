package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"testing"
)

func GenerateUUID() (string, error) {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		return "", err
	}
	//set uuid version and variants acording to RFC 4122
	uuid[6]=(uuid[6] & 0x0f)|0x40 //set version to 4
	uuid[8]=(uuid[8] & 0x3f)|0x80 //set version to 4

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
	
}

func TestUuid(t*testing.T)  {
	id, err:=GenerateUUID()
	if err != nil {
		t.Fatalf("GenerateSellerId failed: %v", err)
	}
	if id=="" {
		t.Log("Succeded\n")		
	}

	// Define a regular expression for UUID v4
	// UUID v4 format: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
	// where y is one of [8, 9, A, B]

	re:=regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
	if !re.MatchString(id){
		t.Errorf("GenerateSellerId returned an invalid UUID: %s", id)
	} else{
		// Log the result if it's valid
		t.Logf("GenerateSellerId succeeded. Generated UUID: %s", id)
	}
}
