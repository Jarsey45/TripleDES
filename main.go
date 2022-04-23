package main

import (
	des "TripleDES/DES"
	"fmt"
)


func main() {
	keys := [3]string{
		"123456789ABCDEF1",
		"6E3272357538782F",
		"292CC794C82AC144",
	}

	var plainText string
	fmt.Print("Enter plaintext to cipher (must be 8 characters long [64-bits]): ")
	fmt.Scan(&plainText)

	cipher := des.Cipher3DES(plainText, keys)
	result := des.Decipher3DES(cipher, keys)
	fmt.Println("\nPLAINTEXT: ", plainText)
	fmt.Println("Triple DES Cipher text: ", cipher)
	fmt.Println("Triple DES Descipher result: ", result)
}
