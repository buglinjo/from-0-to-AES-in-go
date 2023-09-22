package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

var expandedKey [176]byte

func main() {
	/**
	 *
	 *  TODO: add support of 256 bit cipher
	 *
	 */
	message := []byte("This is a message we will encrypt with AES!")
	key := generateKey(128)
	encryptedMessage := encrypt(&message, &key)
	fmt.Println("--------------------------------------")
	fmt.Println(hex.EncodeToString(key[:]))
	fmt.Println(string(message), hex.EncodeToString(message))
	fmt.Println(base64.StdEncoding.EncodeToString(encryptedMessage), hex.EncodeToString(encryptedMessage))
}

func generateKey(bits int) [16]byte {
	//TODO: Generate random key

	return [16]byte{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
}

func encrypt(message *[]byte, key *[16]byte) []byte {
	keyExpansion(key, &expandedKey)

	paddedMessage := addPaddingToMessage(message)
	var encryptedMessage []byte
	for i := 0; i < len(paddedMessage); i += 16 {
		var sliceOfMessage [16]byte
		for j := 0; j < 16; j++ {
			sliceOfMessage[j] = paddedMessage[i+j]
		}

		encrypt16Bytes(&sliceOfMessage, key)

		for j := 0; j < 16; j++ {
			encryptedMessage = append(encryptedMessage, sliceOfMessage[j])
		}
	}

	return encryptedMessage
}

func encrypt16Bytes(state *[16]byte, key *[16]byte) {
	rounds := 10

	addRoundKey(state, key)

	for i := 0; i < rounds; i++ {
		stepsEachRound(state, i)
	}
}

func addPaddingToMessage(message *[]byte) []byte {
	originalLen := len(*message)
	lenOfPaddedMessage := originalLen

	if lenOfPaddedMessage%16 != 0 {
		lenOfPaddedMessage = ((lenOfPaddedMessage / 16) + 1) * 16
	}

	var paddedMessage []byte
	for i := 0; i < lenOfPaddedMessage; i++ {
		if i >= originalLen {
			paddedMessage = append(paddedMessage, 0)
		} else {
			paddedMessage = append(paddedMessage, (*message)[i])
		}
	}

	return paddedMessage
}

func keyExpansion(inputKey *[16]byte, expandedKeys *[176]byte) {
	for i := 0; i < 16; i++ {
		expandedKeys[i] = inputKey[i]
	}

	bytesGenerated := 16
	rconIteration := 1
	var tmp [4]byte

	for bytesGenerated < 176 {
		for i := 0; i < 4; i++ {
			tmp[i] = expandedKeys[i+bytesGenerated-4]
		}
		if bytesGenerated%16 == 0 {
			keyExpansionCore(&tmp, rconIteration)
			rconIteration++
		}

		for i := 0; i < 4; i++ {
			expandedKeys[bytesGenerated] = expandedKeys[bytesGenerated-16] ^ tmp[i]
			bytesGenerated++
		}
	}
}

func keyExpansionCore(in *[4]byte, i int) {
	in[0], in[1], in[2], in[3] = in[1], in[2], in[3], in[0]
	in[0], in[1], in[2], in[3] = Sbox[in[0]], Sbox[in[1]], Sbox[in[2]], Sbox[in[3]]
	in[0] ^= Rcon[i]
}

func addRoundKey(state *[16]byte, roundKey *[16]byte) {
	for i := 0; i < 16; i++ {
		state[i] ^= roundKey[i]
	}
}

func stepsEachRound(state *[16]byte, round int) {
	subBytes(state)
	shiftRows(state)
	if round != 9 {
		mixColumns(state)
	}

	var key [16]byte
	for i := 0; i < 16; i++ {
		key[i] = expandedKey[16*(round+1)+i]
	}

	addRoundKey(state, &key)
}

func subBytes(state *[16]byte) {
	for i := 0; i < 16; i++ {
		state[i] = Sbox[state[i]]
	}
}

func shiftRows(state *[16]byte) {
	for i := 1; i < 4; i++ {
		shiftRowRight(state, i)
	}
}

func shiftRowRight(state *[16]byte, row int) {
	times := row
	for times > 0 {
		tmpFirstElem := state[row]
		for i := 0; i < 3; i++ {
			index := row + (i * 4)
			state[index] = state[index+4]
		}
		state[row+(3*4)] = tmpFirstElem
		times--
	}
}

func mixColumns(state *[16]byte) {
	var tmpState = [16]byte{}
	for i := 0; i < 16; i++ {
		var tmp byte
		stOffset := (i / 4) * 4
		mmOffset := (i * 4) % 16
		for j := 0; j < 4; j++ {
			stIndex := stOffset + j
			mmIndex := mmOffset + j
			switch MulMatrix[mmIndex] {
			case 1:
				tmp ^= state[stIndex]
			case 2:
				tmp ^= Mul2[state[stIndex]]
			case 3:
				tmp ^= Mul3[state[stIndex]]
			}
		}
		tmpState[i] = tmp
	}

	for i := 0; i < 16; i++ {
		state[i] = tmpState[i]
	}
}
