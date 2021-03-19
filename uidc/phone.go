package uidc

// CellphoneNumberMask 手机号敏感处理
func CellphoneNumberMask(cellphoneNumber string) string {
	if len(cellphoneNumber) >= 7 {
		return cellphoneNumber[:3] + "****" + cellphoneNumber[7:]
	} else {
		return "****"
	}
}
