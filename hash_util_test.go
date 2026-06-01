package manor_test

import (
	"fmt"
	"testing"

	manor "github.com/monchickey/manor-go"
)

// go test -v hash_util_test.go

func TestHashUtil(t *testing.T) {
	message := "This is message string."
	hash64Sum := manor.XXHashSum64([]byte(message))
	md5Bytes := manor.MD5Digest([]byte(message))
	md5Hex := manor.MD5HexDigest([]byte(message))
	fmt.Printf("%d, %q, %s\n", hash64Sum, md5Bytes, md5Hex)
	fmt.Println(fmt.Sprintf("%x", md5Bytes) == md5Hex)
}
