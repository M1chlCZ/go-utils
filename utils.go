package utils

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	_ "github.com/bitly/go-simplejson"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"io"
	"log"
	"math"
	mathRand "math/rand"
	"os"
	_ "os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"
)

func InlineIF(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func InlineIFT[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}

func GetENV(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		WrapErrorLog("Error loading .env file")
	}
	return os.Getenv(key)
}

func ReportError(err string, statusCode int) {
	if !strings.Contains(err, "tx_id_UNIQUE") || strings.Contains(err, "Invalid Token, id User") {
		logToFile("")
		logToFile("//// - HTTP ERROR - ////")
		logToFile("HTTP call failed : " + err + "  Status code: " + fmt.Sprintf("%d", statusCode))
		logToFile("////==========////")
		logToFile("")
	}

	// json.NewEncoder(w).Encode(err)
}

func CreateToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	jwtKey := GetENV("JWT_KEY")

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["idUser"] = userId
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := mathRand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func ScheduleFunc(f func(), interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			f()
			//ReportMessage("Scheduled function executed")
		}
	}()
	return ticker
}

func logToFile(message string) {
	f, err := os.OpenFile("api.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v\n", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.Println(message)
	_ = f.Close()
}

func WrapErrorLog(message string) {
	if !strings.Contains(message, "tx_id_UNIQUE") {
		logToFile("//// - ERROR - ////")
		logToFile(message)
		logToFile("////===========////")
		logToFile("")
	}
}

func ReportMessage(message ...string) {
	go func() {
		logToFile(fmt.Sprintf("%s", message))
		logToFile("")
	}()
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func ArrContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GenerateNewPassword(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		if i%8 == 0 && i != 0 {
			b[i] = '-'
		} else {
			mathRand.Seed(time.Now().UnixNano())
			b[i] = letterRunes[mathRand.Intn(len(letterRunes))]
		}
		//b[i] = letterRunes[mathRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateSocialsToken(length int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		mathRand.Seed(time.Now().UnixNano())
		if i%4 == 0 && i != 0 {
			b[i] = '-'
			b[i+1] = letterRunes[mathRand.Intn(len(letterRunes))]
		} else {
			b[i] = letterRunes[mathRand.Intn(len(letterRunes))]
		}
		//b[i] = letterRunes[mathRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ReadFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			WrapErrorLog(err.Error())
			return
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ReadAvatar(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func InTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func Erc20verify(address string, resp func(...string)) bool {
	erc20 := regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	return erc20.MatchString(address)
}

func HashPass(password string) string {
	h := sha256.Sum256([]byte(password))
	str := fmt.Sprintf("%x", h[:])
	return str
}

func RandInt(min int, max int) int {
	return min + mathRand.Intn(max-min)
}
