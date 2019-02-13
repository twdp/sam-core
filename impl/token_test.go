package impl

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestTokenFacadeImpl_EncodeToken(t *testing.T) {
	samClaims := &SamClaims{
		UserName: "fsffjfsifjfsfsdsdsdddsdds",
		StandardClaims: jwt.StandardClaims{
			Id: "1111",
			ExpiresAt: time.Now().Add(10  * time.Minute).UnixNano(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, samClaims)

	fmt.Println(token.SignedString([]byte("hahah_fdfffsdfsfsdjif_")))
}