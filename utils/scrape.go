package utils

import (
	"fmt"
	"strings"
	"time"

	http "github.com/shooyaaa/core/network"
)

//const FundImgUrl = "http://j4.dfcfw.com/charts/pic6/161028.png?v=20210525155946?v=0.9579091525779848"
const FundImgUrl = "http://j4.dfcfw.com/charts/pic6/%v.png?v=%v"
const CryptoCurrencyUrl = "https://production.api.coindesk.com/v2/price/values/%v?start_date=%v&end_date=%v&ohlc=false"

type DataUrl struct {
	Url     string
	Params  map[string]interface{}
	Headers map[string]string
}

func FundRequest(code interface{}) *DataUrl {
	t := time.Now()
	version := fmt.Sprintf(t.Format("200601020405.000"))
	url := fmt.Sprintf(FundImgUrl, code, strings.Replace(version, ".", "", 1))
	headers := make(map[string]string)
	headers[http.HeaderAccept] = strings.Join([]string{http.MimeTypeSvg, http.MimeTypeApng}, ",")
	return &DataUrl{
		Url:     url,
		Headers: headers,
	}
}
