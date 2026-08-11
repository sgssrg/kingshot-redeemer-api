package redeem

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"resty.dev/v3"
)

func Redeem(fid, kid int, code string) (model.RedeemResponse, error) {
	payload := model.RedeemBodyKS{
		Fid: fid,
		Kid: kid,
		Cdk: code,
	}
	return PostSigned(REDEEM_URL, payload)
}

func SignPayload(data model.RedeemBodyKS) map[string]string {
	// Convert struct to map for signing
	m := map[string]string{
		"fid":  strconv.Itoa(data.Fid),
		"kid":  strconv.Itoa(data.Kid),
		"cdk":  data.Cdk,
		"time": data.Time,
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, m[k]))
	}
	str := strings.Join(parts, "&") + SECRET_KEY

	hash := md5.Sum([]byte(str))
	sign := hex.EncodeToString(hash[:])

	m["sign"] = sign
	return m
}

func PostSigned(url string, payload model.RedeemBodyKS) (model.RedeemResponse, error) {
	payload.Time = fmt.Sprintf("%d", time.Now().Unix())

	signed := SignPayload(payload)

	client := resty.New()
	defer client.Close()

	var resp model.RedeemResponse
	var err error
	var res *resty.Response

	// Retry with exponential backoff (max 3 attempts)
	for attempt := 1; attempt <= 3; attempt++ {
		res, err = client.R().
			SetResult(&resp).
			SetFormData(signed).
			Post(url)

		if err != nil {
			slog.Error("HTTP request failed", "error", err, "attempt", attempt)
		} else {
			slog.Info("HTTP Status",
				"status_code", res.StatusCode(),
				"status_text", res.Status(),
				"fid", payload.Fid,
			)

			// If not rate limited (429) or empty response, break immediately
			if res.StatusCode() != http.StatusTooManyRequests &&
				!(resp.ErrCode == 0 && resp.Msg == "") {
				break
			}
		}

		// Exponential backoff: 2s, 4s, 8s
		backoff := time.Duration(float64(time.Second) * float64(attempt) * 1.5)

		slog.Warn("Rate limited or empty response, backing off",
			"fid", payload.Fid,
			"wait", backoff,
		)
		time.Sleep(backoff)
	}

	return resp, err
}
