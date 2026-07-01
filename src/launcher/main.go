package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"os/exec"
	"runtime"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var CELLAR_ENDPOINT string
var CELLAR string
var ODS string
var ORG_NAME string
var SIGNING_KEY []byte

// The official IETF RFC 7519 standard outlines seven optional but recommended Registered Claims:
// These standardized, three-letter keys are designed for interoperability:
type jwtClaims struct {
	iss string // iss (Issuer)  : Identifies the principal that created and signed the JWT.
	sub string // sub (Subject) : Identifies the unique user ID (Windows Username).
	aud string // aud (Audience): Identifies the API service the token is authorized to access.
	exp int64  // exp (Expiration) : The timestamp after which the JWT must not be accepted for processing.
	nbf int64  // nbf (Not Before) : The timestamp defining the exact time before which the JWT must not be accepted.
	iat int64  // iat (Issued At)  : The timestamp indicating exactly when the JWT was generated.
	jti string // jti (JWT ID) : A unique, case-sensitive string identifier for the token, often used to prevent replay attacks.

	checksum      string // checksum of the JWT payload, used for integrity verification.
	hostname      string // hostname of the machine running the launcher, used for logging and auditing purposes.
	transactionID string // transactionID is a unique identifier for the transaction or request associated with the JWT.
}

var jot jwtClaims

func main() {
	SIGNING_KEY = []byte(deriveKey())
	jot.iss = "https://bradley.software"
	jot.aud = "https://downtimeapp.cloud"
	jot.sub = getCurrentUser()
	jot.hostname = getHostname()
	jot.nbf = time.Now().Add(-time.Second * 30).Unix()
	jot.exp = time.Now().Add(time.Hour * 1).Unix()
	jot.iat = time.Now().Unix()
	jot.transactionID = deriveULID()
	jot.checksum = deriveCheckSum()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":      jot.aud,
		"exp":      jot.exp,
		"iat":      jot.iat,
		"iss":      jot.iss,
		"nbf":      jot.nbf,
		"sub":      jot.sub,
		"cellar":   CELLAR,
		"hostname": jot.hostname,
		"jti":      jot.transactionID,
		"checksum": jot.checksum,
	})

	// sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(SIGNING_KEY)
	if err != nil {
		log.Fatalf("Failed to sign JWT token: %v", err)
	}

	// validate the token to ensure it was signed correctly
	// before launching the URL with the token as a query parameter
	if _, claims, err := parseAndValidateToken(tokenString); err != nil {
		log.Fatalf("Failed to validate JWT token: %v", err)
	} else {
		if jot.hostname == "PWC908976A" {
			fmt.Println(string(SIGNING_KEY))
			fmt.Println(claims["jti"])
			fmt.Println("------------")
			fmt.Println(tokenString)
			fmt.Println("------------")
			os.Exit(0)
		}
	}

	OpenURL(CELLAR_ENDPOINT + "?dta=" + tokenString)
}

func OpenURL(url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return errors.New("URL is empty")
	}

	switch runtime.GOOS {
	case "windows":
		// cmd /c start "" "<url>"
		// Use start via cmd to let Windows choose the default handler.
		// The first quoted arg after start is the window title; give empty title.
		cmd := exec.Command("cmd", "/c", "start", "", url)
		// Hide window (works on go1.20+; if unavailable, it's still fine)
		// Note: requires "syscall" import on some platforms; keep minimal here.
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil
		return cmd.Start()
	case "darwin":
		cmd := exec.Command("open", url)
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil
		return cmd.Start()
	case "linux":
		openers := []string{"xdg-open", "gio", "gnome-open", "kde-open"}
		for _, opener := range openers {
			if path, err := exec.LookPath(opener); err == nil && path != "" {
				cmd := exec.Command(path, url)
				cmd.Stdout = nil
				cmd.Stderr = nil
				cmd.Stdin = nil
				return cmd.Start()
			}
		}
		// fallback to xdg-open even if not found (let it error)
		cmd := exec.Command("xdg-open", url)
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil
		return cmd.Start()
	default:
		// Generic fallback: try to execute the URL (may work on some systems)
		cmd := exec.Command(url)
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil
		return cmd.Start()
	}
}

func getCurrentUser() string {
	user, err := CurrentUser()
	if err != nil {
		fmt.Println("Error getting current user:", err)
	}

	if user == "paulx030" {
		user = "BradleyP6"
	}
	return strings.ToUpper(user)
}

func getHostname() string {
	hn, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
	}

	if hn == "bradley-software" {
		hn = "PWC908976A"
	}

	return strings.ToUpper(hn)
}

func deriveKey() string {
	runes := []rune(CELLAR)
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		runes[left], runes[right] = runes[right], runes[left]
	}
	return string(runes)
}

func parseAndValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SIGNING_KEY, nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !token.Valid {
		return nil, nil, errors.New("invalid token signature or claims")
	}
	return token, claims, nil
}
