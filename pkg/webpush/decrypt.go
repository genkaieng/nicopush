package webpush

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"math"

	"golang.org/x/crypto/hkdf"
)

func decodeMessage(msg string) (salt, serverKey, ciphertext []byte, rs int, err error) {
	payload, err := base64.RawURLEncoding.DecodeString(msg)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	rs = (int(payload[16]) << 24) | (int(payload[17]) << 16) | (int(payload[18]) << 8) | int(payload[19])
	keyIdLen := payload[20]
	salt = payload[:16]
	serverKey = payload[21 : 21+keyIdLen]
	ciphertext = payload[21+keyIdLen:]
	return salt, serverKey, ciphertext, rs, nil
}

func computeSharedKey(privKey, serverKey []byte) ([]byte, error) {
	privateKey, err := ecdh.P256().NewPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	remotePub, err := ecdh.P256().NewPublicKey(serverKey)
	if err != nil {
		return nil, err
	}
	return privateKey.ECDH(remotePub)
}

func deriveKeyAndNonce(ikm, auth, pub, senderKey, salt []byte) (key, nonce []byte, err error) {
	authInfo := []byte("WebPush: info\x00")
	authInfo = append(authInfo, pub...)
	authInfo = append(authInfo, senderKey...)

	hash := sha256.New

	hkdfAuth := hkdf.Extract(hash, ikm, auth)
	prkReader := hkdf.Expand(hash, hkdfAuth, authInfo)
	prk := make([]byte, 32)
	if _, err := io.ReadFull(prkReader, prk); err != nil {
		return nil, nil, err
	}

	prk = hkdf.Extract(hash, prk, salt)

	keyInfo := []byte("Content-Encoding: aes128gcm\x00")
	keyReader := hkdf.Expand(hash, prk, keyInfo)
	key = make([]byte, 16)
	if _, err := io.ReadFull(keyReader, key); err != nil {
		return nil, nil, err
	}

	nonceInfo := []byte("Content-Encoding: nonce\x00")
	nonceReader := hkdf.Expand(hash, prk, nonceInfo)
	nonce = make([]byte, 12)
	if _, err := io.ReadFull(nonceReader, nonce); err != nil {
		return nil, nil, err
	}

	return key, nonce, nil
}

func decryptCipherText(data, cek, nonce []byte, rs int) (result []byte) {
	// データをchunkに分割する
	var chunks [][]byte
	for rs < len(data) {
		chunks = append(chunks, data[0:rs:rs])
		data = data[rs:]
	}
	if len(data) > 0 {
		chunks = append(chunks, data)
	}

	// chunkごとに複合して連結
	for i, chunk := range chunks {
		iv := computeNonce(nonce, i)
		decrypted, err := decryptAESGCM(chunk, cek, iv)
		if err != nil {
			log.Println("ERROR", "複合エラー:", err)
			continue
		}
		result = append(result, decrypted...)
	}
	return result
}

func computeNonce(base []byte, seq int) []byte {
	if seq >= int(math.Pow(2, 48)) {
		log.Println("ERROR", "Nonce index is too large  BAD_CRYPTO")
		return nil
	}

	nonce := make([]byte, 12)
	copy(nonce, base[:12])

	for i := 0; i < 6; i++ {
		// インデックスの特定のバイトを取得してnonceに適用します
		// C# の `(byte)(index / System.Math.Pow(256, i)) & 0xff` と同等
		shiftedIndexByte := byte((seq >> (8 * i)) & 0xff)
		nonce[11-i] ^= shiftedIndexByte
	}
	return nonce
}

func decryptAESGCM(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
