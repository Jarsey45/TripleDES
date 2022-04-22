package triple_des

import (
	"fmt"
	"math/rand"
	"math/bits"
	"time"
	"strconv"
)

type Coder struct {
	keys [3]int64;
	permutedKeys [3]string;
	pt [64]int8; // permutation table
	dataBlock string;
	roundShifting [16]int
	cipher string;
}


func stringToBin64(s string) (binString string) {
	strconv.
}

func bin64ToString(s string) (binString string) {
	for _, c := range s {
			binString = fmt.Sprintf("%s%.8b",binString, c)
	}

	c, _ := strconv.ParseInt(binString, 2, 64);
	for  i := 1 ; bits.LeadingZeros64(uint64(c)) > i; i++ {
		binString = "0" + binString;
	}

	return 
}

func permute(key int64, pt [64]int8) string{
	binaryMessage := []byte(strconv.FormatInt(key,2));

	for len(binaryMessage) != 64{
		binaryMessage = append([]byte{49}, binaryMessage...);
	}

	for i := 0; i < len(pt); i++ {
		binaryMessage[i], binaryMessage[pt[i]] = binaryMessage[pt[i]], binaryMessage[i];
	}

	return string(binaryMessage);
}

func permuteString64(key string, pt [64]int8) string{
	binaryMessage := []byte(key);

	for len(binaryMessage) != 64{
		binaryMessage = append([]byte{49}, binaryMessage...);
	}

	for i := 0; i < len(pt); i++ {
		binaryMessage[i], binaryMessage[pt[i]] = binaryMessage[pt[i]], binaryMessage[i];
	}

	return string(binaryMessage);
}

func permuteString32(key string, pt [32]int8) string{
	binaryMessage := []byte(key);

	for len(binaryMessage) < 32{
		binaryMessage = append([]byte{49}, binaryMessage...);
	}

	for i := 0; i < len(pt); i++ {
		binaryMessage[i], binaryMessage[pt[i]] = binaryMessage[pt[i]], binaryMessage[i];
	}

	return string(binaryMessage);
}

func circularShift(s string, places int) (shifted string){ //LEFT shift
	shifted = s;
	for i := 0; i < places; i++{
		tempChar := shifted[0];
		shifted = shifted[1:];
		shifted += string(tempChar);
	}

	return
}

func XOR(first string, second string, size int) (xor string){
	for i := 0; i < size; i++ {
		if(first[i] != second[i]){
			xor += "1";
		}else {
			xor += "0";
		}
	}
	return;
}

func(c *Coder) CreatePermutationTable(){
	//creating permutation table
	for i := 0; i < 64; i++ {
		c.pt[i] = int8(i);
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c.pt), func(i, j int) { c.pt[i], c.pt[j] = c.pt[j], c.pt[i] })
}

func(c *Coder) InitialPermutation(){

	//hardcoded keys
	c.keys = [3]int64{
			0x357638792F423F45,
			0x6E3272357538782F,
			0x292CC794C82AC144,
		};

	//permuting all 3 keys
	for i, key := range(c.keys){
		permutedKey := permute(key, c.pt);
		c.permutedKeys[i] = permutedKey;
		fmt.Println(i, " klucz: ", permutedKey );
	}

	fmt.Println("Message in Binary:", c.dataBlock);
	c.dataBlock = permuteString64(c.dataBlock, c.pt);
	fmt.Println("Permuted message in Binary:", c.dataBlock);

}

func(c *Coder) KeyTransformation(bitkey56 *string, shift int, cp [48]int) (bitkey48 string){
	LBK := (*bitkey56)[:28];
	LBK = circularShift(LBK, shift);
	RBK := (*bitkey56)[28:];
	LBK = circularShift(LBK, shift);
	fmt.Println("Left 28-bit key:", LBK, "Right 28-bit key:", RBK);

	bitkey := LBK + RBK;
	
	for i := 0; i < 48; i ++ {
		bitkey48 += string(bitkey[cp[i] - 1])
	}

	return;
}

func(c *Coder) ExpansionPermutation(bitkey48 string, RPT string) (xor48 string){
	var blockbit4 [8]string;
	for i := 0; i < 8; i++{
		blockbit4[i] = RPT[i*4:(i+1)*4];
	}

	blockbit6 := blockbit4;
	var blockbit48 string;
	for i := 0; i < 8; i++{
		left := i-1;
		right := i+1;
		if i == 0 {
			left = 7;
		} else if i == 7 {
			right = 0;
		}
		blockbit6[i] = blockbit4[left][3:] + blockbit4[i] + blockbit4[right][:1];
		blockbit48 += blockbit6[i];
	}
	fmt.Println("BlockBit48:", blockbit48);

	xor48 = XOR(blockbit48, bitkey48, 48);

	fmt.Println("bitkey48:", bitkey48);
	fmt.Println("XOR48:", xor48);
	return; 
}

//this function is whatever you can imagine as long as it returns 32-bit block
func(c *Coder) SBox(xor48 string) (blockbit48 string){
	for i := 0; i < 8; i++{
		blockbit48 += xor48[i*6:i*6+4];
	}
	fmt.Println("Blockbit48 (S-BOXES):", blockbit48);
	return;
}

func(c *Coder) PBox(sbox32 string, pBoxTable [32]int8) (pbox32 string) {
	pbox32 = permuteString32(sbox32, pBoxTable);
	fmt.Println("BlockBit32 (P-BOX):", pbox32);
	return;
}

func(c *Coder) xorAndSwap(LPT *string, pbox32 *string, RPT *string){
	//xor
	newRPT := XOR(*LPT, *pbox32, 32);

	//swap
	*LPT = *RPT;
	*RPT = newRPT;
}

func(c *Coder) Rounds(){
	bitkey56 := c.permutedKeys[0][:56];
	LPT := c.dataBlock[:32];
	RPT := c.dataBlock[32:];
	fmt.Println("LPT:", LPT);
	fmt.Println("LPT:", RPT);

	compressionPermutation := [48]int{
		14,17,11,24,1,5,3,28,15,6,21,10,
		23,19,12,4,26,8,16,7,27,20,13,2,
		41,52,31,37,47,55,30,40,51,45,33,48,
		44,49,39,56,34,53,46,42,50,36,29,32,
	}; 

	pBoxTable := [32]int8{ 
		15, 6, 19, 20,
		28, 11, 27, 16,
		0, 14, 22, 25,
		4, 17, 30, 9,
		1, 7, 23, 13,
		31, 26, 2, 8,
		18, 12, 29, 5,
		21, 10, 3, 24 };

	//16 rund
	for i := 0; i < 16; i++ {
		//Key transformation
		bitkey48 :=	c.KeyTransformation(&bitkey56, c.roundShifting[i], compressionPermutation);

		//Expansion Permutation
		xor48 := c.ExpansionPermutation(bitkey48, RPT);

		//S-Box Substitution
		sbox32 := c.SBox(xor48);

		//P-Box Permutation
		pbox32 := c.PBox(sbox32, pBoxTable);

		//XOR and Swap
		c.xorAndSwap(&LPT, &pbox32, &RPT);
	}

	c.dataBlock = LPT + RPT;
}

func(c *Coder) FinalPermutation(){
	finalPermTable := [64]int8{
		39, 7, 47, 15, 55, 23, 63, 31,
		38, 6, 46, 14, 54, 22, 62, 30,
		37, 5, 45, 13, 53, 21, 61, 29,
		36, 4, 44, 12, 52, 20, 60, 28,
		35, 3, 43, 11, 51, 19, 59, 27,
		34, 2, 42, 10, 50, 18, 58, 26,
		33, 1, 41, 9, 49, 17, 57, 25,
		32, 0, 40, 8, 48, 16, 56, 24};
	
	c.cipher = permuteString64(c.dataBlock, finalPermTable);

}

func(c *Coder) EncodeDES(data string){
	c.dataBlock = stringToBin64("tata");
	c.roundShifting = [16]int{1,1,2,2,2,2,2,2,1,2,2,2,2,2,2,1};

	//Krok 1: permutacja początkowa
	c.InitialPermutation();
	//Krok 2: rundy
	c.Rounds();
	//KROK 3: permutacja finalna
	c.FinalPermutation();

	fmt.Println("CIPHER (string): ", c.cipher);
}

func(c *Coder) Decode(data string){

}
