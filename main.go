package main

import (
	des "TripleDES/DES"
	"encoding/hex"
	"fmt"
	"strings"
)

// var coder cipher.Coder;
var coder des.DES
func main() {
	keys := [3]string{
		"123456789ABCDEF1",
		"6E3272357538782F",
		"292CC794C82AC144",
	}

	plainText := "8bitow!!"

	cipher := cipher3DES(plainText, keys);
	result := decipher3DES(cipher, keys);
	fmt.Println("Triple DES Cipher text: ", cipher);
	fmt.Println("Triple DES Descipher result: ", result);
}

func cipher3DES(plainText string, keys [3]string) (cipher string){
	plainTextHex := strings.ToUpper(hex.EncodeToString([]byte(plainText)))

	//CIPHERTEXT
	//Encode(k1)
	encode1 := coder.EncodeDES(plainTextHex, keys[0]);
	//Decode(k2)
	encode2 := coder.DecodeDES(keys[1], encode1);
	//Encode(k3)
	cipher = coder.EncodeDES(encode2, keys[2]);
	return
}

func decipher3DES(cipher string, keys [3]string) (plainText string){

	//DECIPHERTEXT
	//Decode(k3)
	decode1 := coder.DecodeDES(keys[2], cipher);
	//Encode(k2)
	decode2 := coder.EncodeDES(decode1, keys[1]);
	//Decode(k1)
	decode3 := coder.DecodeDES(keys[0], decode2);

	ascii ,_ := hex.DecodeString(decode3);
	plainText = string(ascii);
	return
}