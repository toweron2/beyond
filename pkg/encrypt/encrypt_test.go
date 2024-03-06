package encrypt

import (
    "testing"
)

func TestEncryptMobile(t *testing.T) {
    mobile := "13800138000"
    encryptedMobile, err := EncMobile(mobile)
    if err != nil {
        t.Fatal(err)
    }
    decrypteMobile, err := DecMobile(encryptedMobile)
    if err != nil {
        t.Fatal(err)
    }
    if mobile != decrypteMobile {
        t.Fatalf("expected %s, but got %s", mobile, decrypteMobile)
    }
}
