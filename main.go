package main

import (
	des "TripleDES/DES"
	"fmt"
)

// var coder cipher.Coder;
func main() {
	keys := [3]string{
		"123456789ABCDEF1",
		"6E3272357538782F",
		"292CC794C82AC144",
	}

	plainText := "8bitow!!"

	cipher := des.Cipher3DES(plainText, keys);
	result := des.Decipher3DES(cipher, keys);
	fmt.Println("PLAINTEXT: ", plainText);
	fmt.Println("Triple DES Cipher text: ", cipher);
	fmt.Println("Triple DES Descipher result: ", result);
}

