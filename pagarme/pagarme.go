package pagarme

import (
	"fmt"
	"github.com/leekchan/accounting"
	"strconv"
	"strings"
)

func FormatAmount(amount float64) int64 {
	ac := accounting.Accounting{Symbol: "", Precision: 2, Thousand: "", Decimal: ""}
	text := strings.Replace(ac.FormatMoney(amount), ",", "", -1)
	text = strings.Replace(text, ".", "", -1)
	val, _ := strconv.Atoi(text)
	return int64(val)
}

func FormatAmountToFloat(amount int64) float64 {
	unformateed := accounting.UnformatNumber(fmt.Sprintf("%v", amount), 2, "BRL")
	val, _ := strconv.ParseFloat(unformateed, 64)
	return val
}
