package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secretpassword")

var KeyPath = "E:\\G-projects\\Keys\\"
var BitSize = 4096

func KeyWrite(c *gin.Context) {
	//To delte prevois files
	files, err := filepath.Glob(filepath.Join("E:\\G-projects\\Keys", "*"))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			panic(err)
		}
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, BitSize)
	if err != nil {
		panic(err)
	}
	// Private Key
	pvtraw := x509.MarshalPKCS1PrivateKey(privateKey)
	pvtblock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pvtraw,
	}
	file, err := os.Create(KeyPath + "private.pfx")
	if err != nil {
		panic(err)
	}
	err = pem.Encode(file, pvtblock)
	if err != nil {
		panic(err)
	}
	// Public key
	pubraw := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubblock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubraw,
	}
	pubfile, err := os.Create(KeyPath + "public.pem")
	if err != nil {
		panic(err)
	}
	err = pem.Encode(pubfile, pubblock)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": "Keys generated successfully"})
}

func Temp(pass string) string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey
	encryptedBytes, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		&publicKey,
		[]byte("super secret message"),
		nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("encrypted bytes: ", encryptedBytes)
	decryptedBytes, err := privateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA512})
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted message: ", string(decryptedBytes))
	return string(decryptedBytes)
}

func EncryptPassword(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func VerifyPassword(hashPass []byte, pass []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashPass, pass)
	return err == nil
}

func CreateToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	})
	//token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
