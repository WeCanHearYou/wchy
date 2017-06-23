package jwt_test

import (
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/jwt"
	. "github.com/onsi/gomega"
)

func TestJWT_Encode(t *testing.T) {
	RegisterTestingT(t)

	claims := &models.FiderClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, err := jwt.Encode(claims)
	Expect(token).NotTo(BeNil())
	Expect(err).To(BeNil())
}

func TestJWT_Decode(t *testing.T) {
	RegisterTestingT(t)

	claims := &models.FiderClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.DecodeFiderClaims(token)
	Expect(err).To(BeNil())
	Expect(decoded.UserID).To(Equal(claims.UserID))
	Expect(decoded.UserName).To(Equal(claims.UserName))
	Expect(decoded.UserEmail).To(Equal(claims.UserEmail))
}

func TestJWT_DecodeOAuthClaims(t *testing.T) {
	RegisterTestingT(t)

	claims := &models.OAuthClaims{
		OAuthID:       "2",
		OAuthEmail:    "jon.snow@got.com",
		OAuthName:     "Jon Snow",
		OAuthProvider: "facebook",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.DecodeOAuthClaims(token)
	Expect(err).To(BeNil())
	Expect(decoded.OAuthID).To(Equal(claims.OAuthID))
	Expect(decoded.OAuthEmail).To(Equal(claims.OAuthEmail))
	Expect(decoded.OAuthName).To(Equal(claims.OAuthName))
	Expect(decoded.OAuthProvider).To(Equal(claims.OAuthProvider))
}

func TestJWT_DecodeChangedToken(t *testing.T) {
	RegisterTestingT(t)

	claims := &models.FiderClaims{
		UserID:    424,
		UserName:  "Jon Snow",
		UserEmail: "jon.snow@got.com",
	}

	token, _ := jwt.Encode(claims)

	decoded, err := jwt.DecodeFiderClaims(token + "foo")
	Expect(err).ToNot(BeNil())
	Expect(decoded).To(BeNil())
}
