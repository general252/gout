package uidc

import "log"

func ExampleIdCardNumberMask() {
	log.Println(IdCardNumberMask("340102199003070037"))
	log.Println(IdCardNumberMask("34010219900307429X"))
	// output:
}

func ExampleIdCardNumberCheck() {
	log.Println(IdCardNumberCheck("110101199003078291"))
	// output:
}

func ExampleGetBirthdayFromIdCardNumber() {
	log.Println(GetBirthdayFromIdCardNumber("110101199009218697"))
	// output:
}
