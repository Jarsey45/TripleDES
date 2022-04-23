package des

import (
	"fmt"
	"strconv"
)

type DES struct {
	key         string
	keyBinary  	string
	roundKeys    [16]string
	roundKeysHex [16]string
	binaryText   string
	cipherText   string
	stage int
}

func reverseRoundKey(arr *[16]string) {
	for i, j := 0, 15; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func hex2bin(s string) (bin string) {
	mp := map[string]string{
		"0": "0000",
		"1": "0001",
		"2": "0010",
		"3": "0011",
		"4": "0100",
		"5": "0101",
		"6": "0110",
		"7": "0111",
		"8": "1000",
		"9": "1001",
		"A": "1010",
		"B": "1011",
		"C": "1100",
		"D": "1101",
		"E": "1110",
		"F": "1111",
	}
	for i := 0; i < len(s); i++ {
		bin += mp[string(s[i])]
	}
	return
}

func bin2hex(s string) (bin string) {
	mp := map[string]string{
		"0000": "0",
		"0001": "1",
		"0010": "2",
		"0011": "3",
		"0100": "4",
		"0101": "5",
		"0110": "6",
		"0111": "7",
		"1000": "8",
		"1001": "9",
		"1010": "A",
		"1011": "B",
		"1100": "C",
		"1101": "D",
		"1110": "E",
		"1111": "F",
	}
	for i := 0; i < len(s)/4; i++ {
		bin += mp[string(s[i*4:(i+1)*4])]
	}
	return
}

func dec2bin(s int64) (bin string) {
	bin = fmt.Sprintf("%04b", s);
	return;
}

func bin2dec(s string) (dec int64) {
	dec, _ = strconv.ParseInt(s, 2, 64);
	return 
}

func permute(key string, pt []int8, size int) string {
	permutation := "";
	for i := 0; i < size; i++ {
		permutation += string(key[pt[i] - 1])
	}
	return permutation;
}

func shift(s string, places int8) (shifted string) { //LEFT shift
	shifted = s
	for i := int8(0); i < places; i++ {
		tempChar := shifted[0]
		shifted = shifted[1:]
		shifted += string(tempChar)
	}

	return
}

func XOR(first string, second string, size int) (xor string) {
	for i := 0; i < size; i++ {
		if first[i] != second[i] {
			xor += "1"
		} else {
			xor += "0"
		}
	}
	return
}

func swap(str1, str2 *string) {
	*str1, *str2 = *str2, *str1
}

func (c *DES) GenerateKeys() {
	//convert all Keys to binary
	c.keyBinary = hex2bin(c.key);

	keyPermutationTable := []int8{
		57, 49, 41, 33, 25, 17, 9,
		1, 58, 50, 42, 34, 26, 18,
		10, 2, 59, 51, 43, 35, 27,
		19, 11, 3, 60, 52, 44, 36,
		63, 55, 47, 39, 31, 23, 15,
		7, 62, 54, 46, 38, 30, 22,
		14, 6, 61, 53, 45, 37, 29,
		21, 13, 5, 28, 20, 12, 4,
	}

	//permute all keys
	c.keyBinary = permute(c.keyBinary, keyPermutationTable, 56);

	shiftTable := [16]int8{
		1, 1, 2, 2,
		2, 2, 2, 2,
		1, 2, 2, 2,
		2, 2, 2, 1,
	}

	keyCompressionTable := []int8{
		14, 17, 11, 24, 1, 5,
		3, 28, 15, 6, 21, 10,
		23, 19, 12, 4, 26, 8,
		16, 7, 27, 20, 13, 2,
		41, 52, 31, 37, 47, 55,
		30, 40, 51, 45, 33, 48,
		44, 49, 39, 56, 34, 53,
		46, 42, 50, 36, 29, 32,
	}


	leftKey := c.keyBinary[:28]
	rightKey := c.keyBinary[28:]


	for i := 0; i < 16; i++ {
		leftKey = shift(leftKey, shiftTable[i])
		rightKey = shift(rightKey, shiftTable[i])
		
		combineBitkey := leftKey + rightKey

		roundKey := permute(combineBitkey, keyCompressionTable, 48)
		
		c.roundKeys[i] = roundKey
		c.roundKeysHex[i] = bin2hex(roundKey)
		
		
	}
}

func (c *DES) InitialPermutation() {
	initialPermutationTable := []int8{
		58, 50, 42, 34, 26, 18, 10, 2,
		60, 52, 44, 36, 28, 20, 12, 4,
		62, 54, 46, 38, 30, 22, 14, 6,
		64, 56, 48, 40, 32, 24, 16, 8,
		57, 49, 41, 33, 25, 17, 9, 1,
		59, 51, 43, 35, 27, 19, 11, 3,
		61, 53, 45, 37, 29, 21, 13, 5,
		63, 55, 47, 39, 31, 23, 15, 7,
	}

	c.binaryText = permute(c.binaryText, initialPermutationTable, 64)
}

func (c *DES) ExpansionPermutation(rightText32 string, roundKey string) (xor48 string) {

	//expanding 32-bit key to 48-bit
	var blockbit4 [8]string
	for i := 0; i < 8; i++ {
		blockbit4[i] = rightText32[i*4 : (i+1)*4]
	}

	blockbit6 := blockbit4
	var blockbit48 string
	for i := 0; i < 8; i++ {
		left := i - 1
		right := i + 1
		if i == 0 {
			left = 7
		} else if i == 7 {
			right = 0
		}
		blockbit6[i] = blockbit4[left][3:] + blockbit4[i] + blockbit4[right][:1]
		blockbit48 += blockbit6[i]
	}

	//XOR-ed key with roundKey[i]
	xor48 = XOR(blockbit48, roundKey, 48)

	return
}

func (c *DES) SBox(xorKey48 string) (sbox string){
	sboxTable := [][][]int8{{{14, 4, 13, 1, 2, 15, 11, 8, 3, 10, 6, 12, 5, 9, 0, 7},
          { 0, 15, 7, 4, 14, 2, 13, 1, 10, 6, 12, 11, 9, 5, 3, 8},
          { 4, 1, 14, 8, 13, 6, 2, 11, 15, 12, 9, 7, 3, 10, 5, 0},
          {15, 12, 8, 2, 4, 9, 1, 7, 5, 11, 3, 14, 10, 0, 6, 13 }},
            
         {{15, 1, 8, 14, 6, 11, 3, 4, 9, 7, 2, 13, 12, 0, 5, 10},
            {3, 13, 4, 7, 15, 2, 8, 14, 12, 0, 1, 10, 6, 9, 11, 5},
            {0, 14, 7, 11, 10, 4, 13, 1, 5, 8, 12, 6, 9, 3, 2, 15},
           {13, 8, 10, 1, 3, 15, 4, 2, 11, 6, 7, 12, 0, 5, 14, 9 }},
   
         { {10, 0, 9, 14, 6, 3, 15, 5, 1, 13, 12, 7, 11, 4, 2, 8},
           {13, 7, 0, 9, 3, 4, 6, 10, 2, 8, 5, 14, 12, 11, 15, 1},
           {13, 6, 4, 9, 8, 15, 3, 0, 11, 1, 2, 12, 5, 10, 14, 7},
            {1, 10, 13, 0, 6, 9, 8, 7, 4, 15, 14, 3, 11, 5, 2, 12 }},
       
          { {7, 13, 14, 3, 0, 6, 9, 10, 1, 2, 8, 5, 11, 12, 4, 15},
           {13, 8, 11, 5, 6, 15, 0, 3, 4, 7, 2, 12, 1, 10, 14, 9},
           {10, 6, 9, 0, 12, 11, 7, 13, 15, 1, 3, 14, 5, 2, 8, 4},
            {3, 15, 0, 6, 10, 1, 13, 8, 9, 4, 5, 11, 12, 7, 2, 14} },
        
          { {2, 12, 4, 1, 7, 10, 11, 6, 8, 5, 3, 15, 13, 0, 14, 9},
           {14, 11, 2, 12, 4, 7, 13, 1, 5, 0, 15, 10, 3, 9, 8, 6},
            {4, 2, 1, 11, 10, 13, 7, 8, 15, 9, 12, 5, 6, 3, 0, 14},
           {11, 8, 12, 7, 1, 14, 2, 13, 6, 15, 0, 9, 10, 4, 5, 3 }},
       
         { {12, 1, 10, 15, 9, 2, 6, 8, 0, 13, 3, 4, 14, 7, 5, 11},
           {10, 15, 4, 2, 7, 12, 9, 5, 6, 1, 13, 14, 0, 11, 3, 8},
            {9, 14, 15, 5, 2, 8, 12, 3, 7, 0, 4, 10, 1, 13, 11, 6},
            {4, 3, 2, 12, 9, 5, 15, 10, 11, 14, 1, 7, 6, 0, 8, 13} },
         
          { {4, 11, 2, 14, 15, 0, 8, 13, 3, 12, 9, 7, 5, 10, 6, 1},
           {13, 0, 11, 7, 4, 9, 1, 10, 14, 3, 5, 12, 2, 15, 8, 6},
            {1, 4, 11, 13, 12, 3, 7, 14, 10, 15, 6, 8, 0, 5, 9, 2},
            {6, 11, 13, 8, 1, 4, 10, 7, 9, 5, 0, 15, 14, 2, 3, 12} },
        
         { {13, 2, 8, 4, 6, 15, 11, 1, 10, 9, 3, 14, 5, 0, 12, 7},
            {1, 15, 13, 8, 10, 3, 7, 4, 12, 5, 6, 11, 0, 14, 9, 2},
            {7, 11, 4, 1, 9, 12, 14, 2, 0, 6, 10, 13, 15, 3, 5, 8},
            {2, 1, 14, 7, 4, 10, 8, 13, 15, 12, 9, 0, 3, 5, 6, 11} } }

	for i := 0; i < 8; i++ {
			rindex := string(xorKey48[i * 6]) + string(xorKey48[i * 6 + 5]);
			row := bin2dec(rindex);
			cindex := string(xorKey48[i * 6 + 1]) + string(xorKey48[i * 6 + 2]) + string(xorKey48[i * 6 + 3]) + string(xorKey48[i * 6 + 4]);
			col := bin2dec(cindex);
			val := sboxTable[i][row][col];
			sbox += dec2bin(int64(val));
	}
		
	return;
}

func (c *DES) Rounds() (combine string){
	sboxPermutationTable := []int8{
		16,  7, 20, 21,
		29, 12, 28, 17,
		1, 15, 23, 26,
		5, 18, 31, 10,
		2,  8, 24, 14,
		32, 27,  3,  9,
		19, 13, 30,  6,
		22, 11,  4, 25,
	}

	leftText32 := c.binaryText[:32]
	rightText32 := c.binaryText[32:]


	for i := 0; i < 16; i++ {
		//D-BOX: Expanding the 32 bits data into 48 bits
		rightText48 := c.ExpansionPermutation(rightText32, c.roundKeys[i])
		
		//S-BOX: substituting the value from s-box table by calculating row and column
		sbox32 := c.SBox(rightText48);
		
		//Straight D-BOX
		sbox32 = permute(sbox32, sboxPermutationTable,32);
		
		//XOR left and sbox32
		result := XOR(leftText32, sbox32, 32);
		leftText32 = result;
		
		if(i != 15){

			swap(&leftText32, &rightText32);
		}

		fmt.Println("Round ", i + 1, " ", bin2hex(leftText32), " ", bin2hex(rightText32), " ", c.roundKeysHex[i])
	}
	combine = leftText32 + rightText32;
	return
}

func(c *DES) FinalPermutation(combination string){
	finalPermTable := []int8{
		40, 8, 48, 16, 56, 24, 64, 32,
		39, 7, 47, 15, 55, 23, 63, 31,
		38, 6, 46, 14, 54, 22, 62, 30,
		37, 5, 45, 13, 53, 21, 61, 29,
		36, 4, 44, 12, 52, 20, 60, 28,
		35, 3, 43, 11, 51, 19, 59, 27,
		34, 2, 42, 10, 50, 18, 58, 26,
		33, 1, 41, 9, 49, 17, 57, 25};
	
	c.cipherText = bin2hex(permute(combination, finalPermTable, 64));

}

func (c *DES) Encrypt() {
	//Initial Permutation
	c.InitialPermutation()

	//Rounds
	combination := c.Rounds();

	c.FinalPermutation(combination);

}

func (c *DES) EncodeDES(data string, key string) string {
	fmt.Println("-----ENCODING-----");
	//change key
	c.key = key;

	c.GenerateKeys()

	//convert data to binary
	c.binaryText = hex2bin(data)
	fmt.Println("TEXT (HEX): ", data);


	//ENCRYPT
	c.Encrypt()

	fmt.Println("CIPHER TEXT", c.cipherText)
	return c.cipherText
}

func (c *DES) DecodeDES(key string, cipher string) string {
	fmt.Println("-----DECODING-----");
	//change key
	c.key = key;

	c.GenerateKeys()

	//reverse roundKeys
	reverseRoundKey(&c.roundKeys);
	reverseRoundKey(&c.roundKeysHex);

	c.binaryText = hex2bin(cipher);
	fmt.Println("TEXT (HEX): ", cipher);

	// II: ENCRYPT
	c.Encrypt();

	fmt.Println("DECIPHERED TEXT (HEX):", c.cipherText)
	return c.cipherText
}
