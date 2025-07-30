package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

func (s *Data) buildSecretIndex() {
	if s.secretIndex == nil {
		s.secretIndex = make(map[string]int)
	} else {
		for k := range s.secretIndex {
			delete(s.secretIndex, k)
		}
	}

	for i, secret := range s.Secrets {
		s.secretIndex[secret.Name] = i
	}
}

func (s *Data) ensureSecretIndex() {
	if s.secretIndex == nil || len(s.secretIndex) != len(s.Secrets) {
		s.buildSecretIndex()
	}
}

func (s *Data) addToSecretIndex(secretName string, index int) {
	if s.secretIndex == nil {
		s.secretIndex = make(map[string]int)
	}
	s.secretIndex[secretName] = index
}

func (s *Data) removeFromSecretIndex(secretName string, removedIndex int) {
	if s.secretIndex == nil {
		return
	}

	delete(s.secretIndex, secretName)

	for name, idx := range s.secretIndex {
		if idx > removedIndex {
			s.secretIndex[name] = idx - 1
		}
	}
}

func (s *Data) updateSecretIndex(oldName, newName string, index int) {
	if s.secretIndex == nil {
		s.secretIndex = make(map[string]int)
	}

	delete(s.secretIndex, oldName)
	s.secretIndex[newName] = index
}

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

func (s *Data) checkIfSecretExists(name string) (bool, int) {
	s.ensureSecretIndex()

	if index, exists := s.secretIndex[name]; exists {
		return true, index
	}
	return false, -1
}

func (s *Data) fuzzySearchSecret(name string) (int, error) {
	if len(s.Secrets) == 0 {
		return -1, errors.New("secret was not found")
	}

	idToIndex := make(map[string]int, len(s.Secrets))
	secretIds := make([]string, len(s.Secrets))

	for i, secret := range s.Secrets {
		secretIds[i] = secret.Name
		idToIndex[secret.Name] = i
	}

	results := fuzzy.RankFind(name, secretIds)

	if len(results) == 0 {
		return -1, errors.New("secret was not found")
	}

	bestMatch := results[0].Target
	return idToIndex[bestMatch], nil
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
