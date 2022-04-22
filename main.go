package main

import (
	cipher "TripleDES/3DES"
)

var coder cipher.Coder;
func main() {
	coder.CreatePermutationTable();
	coder.EncodeDES("8bitow!!"); //8 znaków wiadomosci
}