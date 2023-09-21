package main

func main() {
	//a := [4]byte{1, 2, 3, 4}
	//keyExpansionCore(&a, 1)
	//fmt.Println(a)

	/**
	 *
	 *  TODO: add support of 256 bit cipher
	 *
	 */
	message := []byte("This is the message we will encrypt with AES")
	key := generateKey(128)
	encrypt(&message, &key)
}

func generateKey(bits int) []byte {
	//TODO: Generate random key

	return []byte{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
}

func encrypt(message *[]byte, key *[]byte) {
	rounds := 10
	var state [16]byte

	for i := 0; i < 16; i++ {
		state[i] = (*message)[i]
	}

	keyExpansion()
	addRoundKey(&state, key)

	for i := 0; i < rounds; i++ {
		stepsEachRound(&state, key, i == rounds-1)
	}
}

func keyExpansion() {
	//TODO add key expansion functionality
}

func keyExpansionCore(in *[4]byte, i byte) {
	in[0], in[1], in[2], in[3] = in[1], in[2], in[3], in[0]
	in[0], in[1], in[2], in[3] = Sbox[in[0]], Sbox[in[1]], Sbox[in[2]], Sbox[in[3]]
}

func addRoundKey(state *[16]byte, roundKey *[]byte) {
	for i := 0; i < 16; i++ {
		state[i] ^= (*roundKey)[i]
	}
}

func stepsEachRound(state *[16]byte, key *[]byte, isFinalRound bool) {
	subBytes(state)
	shiftRows(state)

	if !isFinalRound {
		mixColumns(state)
	}

	addRoundKey(state, key)
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
		stOffset := (i * 4) % 16
		mmOffset := (i / 4) * 4
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
