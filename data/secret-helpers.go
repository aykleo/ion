package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *Data) extractArgs(args []string, generateSaltIfMissing bool) (string, string, string, []string, error) {
	name, value, err := extractNameAndValue(args)
	if err != nil {
		return "", "", "", nil, err
	}

	hasSalt, salt := detectSalt(args)
	hasTags, tags := detectTags(args)

	if !hasSalt && generateSaltIfMissing {
		salt = generateSalt()
	}
	tgs := []string{}
	if hasTags {
		tgs = tags
	}

	return name, value, salt, tgs, nil
}

func generateSalt() string {
	mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", mathrand.Intn(1000000))
}

func encrypt(salt string, value string) string {
	key := sha256.Sum256([]byte(salt))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return ""
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return ""
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)

	return hex.EncodeToString(ciphertext)
}

func decrypt(salt string, encryptedValue string) string {
	key := sha256.Sum256([]byte(salt))

	ciphertext, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return ""
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return ""
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return ""
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return ""
	}

	return string(plaintext)
}

func (s *Data) checkIfSecretExists(name string, path string) (bool, int) {
	dataPath := filepath.Join(path, "data.json")
	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		return false, -1
	}
	var data Data
	if err := json.Unmarshal(fileData, &data); err != nil {
		return false, -1
	}

	for i, secret := range data.Secrets {
		if secret.ID == name {
			return true, i
		}
	}
	return false, -1
}

func detectSalt(args []string) (bool, string) {
	for i, arg := range args {
		if arg == "-s" || arg == "--salt" {
			if i+1 < len(args) {
				return true, args[i+1]
			}
			return false, ""
		}
	}
	return false, ""
}

func detectTags(args []string) (bool, []string) {
	tags := []string{}
	for i, arg := range args {
		if arg == "-t" || arg == "--tags" {
			for j := i + 1; j < len(args); j++ {
				if args[j] == "-s" || args[j] == "--salt" {
					break
				}
				if j >= len(args)-2 {
					break
				}
				tags = append(tags, strings.ToUpper(args[j]))
			}
			break
		}
	}
	return len(tags) > 0, tags
}

func extractNameAndValue(args []string) (string, string, error) {
	usedIndices := make(map[int]bool)

	for i, arg := range args {
		if arg == "-t" || arg == "--tags" {
			usedIndices[i] = true
			for j := i + 1; j < len(args); j++ {
				if args[j] == "-s" || args[j] == "--salt" {
					break
				}
				if j >= len(args)-2 {
					break
				}
				if len(args[j]) > 0 && args[j][0] == '-' {
					break
				}
				usedIndices[j] = true
			}
		}
		if arg == "-s" || arg == "--salt" {
			usedIndices[i] = true
			if i+1 < len(args) {
				usedIndices[i+1] = true
			}
		}
		if arg == "-k" || arg == "--keep-value" {
			usedIndices[i] = true
		}
		if arg == "-n" || arg == "--name" {
			usedIndices[i] = true
		}
	}

	var remaining []string
	for i, arg := range args {
		if !usedIndices[i] {
			remaining = append(remaining, arg)
		}
	}

	if len(remaining) != 2 {
		return "", "", errors.New("expected exactly name and value")
	}

	return remaining[0], remaining[1], nil
}

// func detectKeepValue(args []string) (bool, int) {
// 	for i, arg := range args {
// 		if arg == "-k" || arg == "--keep-value" {
// 			return true, i
// 		}
// 	}
// 	return false, -1
// }

// func detectName(args []string) (bool, int) {
// 	for i, arg := range args {
// 		if arg == "-n" || arg == "--name" {
// 			return true, i
// 		}
// 	}
// 	return false, -1
// }
