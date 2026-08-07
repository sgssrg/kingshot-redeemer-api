package redeem

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"resty.dev/v3"
)

const (
	REDEEM_URL = "https://kingshot-giftcode.centurygame.com/api/gift_code"
	SECRET_KEY = "mN4!pQs6JrYwV9"
)

// RedeemWithDB redeems a gift code for all players in the database.
// @Summary Redeem gift code for all players
// @Description Loads all players from the database and redeems the configured gift code for each one, streaming per-player results as Server-Sent Events (event: redeem-res) followed by a summary (event: redeem-fin). Rate limited.
// @Tags Redeem
// @Produce text/event-stream
// @Success 200 {string} string "SSE stream of redeem-res and redeem-fin events"
// @Router /redeem/db [post]
func RedeemWithDB(c *echo.Context) error {
	reqContext := c.Request().Context()
	res := c.Response()

	// 1. Set required SSE headers
	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("X-Accel-Buffering", "no") // Prevents Nginx/reverse proxies from buffering the stream

	// Ensure the underlying ResponseWriter supports HTTP flushing
	flusher, ok := res.(http.Flusher)
	if !ok {
		slog.Error("Streaming unsupported: response writer does not implement http.Flusher")
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
	}

	giftCode := c.QueryParam("code")
	if giftCode == "" {
		giftCode = "KS0803"
	}
	dbApp := db.InitDB()
	fmt.Println(os.Getenv("DB_URI"))
	defer dbApp.Pool.Close()

	ctx := context.Background()
	players, err := dbApp.Queries.GetAllPlayers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	chunks := chunkArray(players, 5)
	var playerMetaRedeemResponse model.PlayerMetaRedeemResponse
	var finalRes []model.RedeemSSE
	for _, batch := range chunks {
		for _, p := range batch {
			// Check if the client disconnected before doing work
			select {
			case <-reqContext.Done():
				slog.Info("Client disconnected, stopping process")
				return nil
			default:
			}

			resp, err := redeem(int(p.Pid), int(p.Kid), giftCode)
			redeemEvent := model.RedeemSSE{
				Fid:     strconv.Itoa(int(p.Pid)),
				Code:    giftCode,
				Result:  resp,
				Success: err == nil,
				Time:    time.Now().Format(time.RFC3339),
			}
			slog.Info("resp", "data", resp)
			finalRes = append(finalRes, redeemEvent)
			if err == nil && resp.ErrCode == 0 && resp.Msg == "" {
				slog.Warn("Received empty response, retrying in 5 seconds...", "fid", int(p.Pid))
				sleep(5)
				resp, err = redeem(int(p.Pid), 1420, giftCode)
			}
			switch redeemEvent.Result.ErrCode {
			case 20000:
				playerMetaRedeemResponse.Redeemed++
			case 40008:
				playerMetaRedeemResponse.Manual++
			case 40020:
				playerMetaRedeemResponse.PlayerDead++
			case 40005, 40007, 40014:
				playerMetaRedeemResponse.CodeExpired++
			}

			if err != nil {
				slog.Error("Error in Redeem", "error", err)
			}

			// 2. Marshal payload to JSON
			jsonData, err := json.Marshal(redeemEvent)
			if err != nil {
				slog.Error("Failed to marshal event", "error", err)
				continue
			}

			// 3. Format as standard SSE payload: "data: <json>\n\n"
			if _, err := fmt.Fprintf(res, "event: redeem-res\ndata: %s\n\n", jsonData); err != nil {
				return err
			}

			// 4. Immediately flush buffer to the client
			flusher.Flush()
			// sleep(1)
		}
		// Delay between batches
		for i := 0; i < 4; i++ {
			sleep(5)
			// Send an SSE comment to keep the TCP pipe active
			fmt.Fprintf(res, ": keep-alive\n\n")
			flusher.Flush()
		}
		// sleep(1)
	}
	metaJsonData, err := json.Marshal(playerMetaRedeemResponse)
	if err != nil {
		slog.Error("Failed to marshal event", "error", err)
		return nil
	}

	// FIX 2: Fixed duplicate "data: data:" -> "data: %s\n\n"
	fmt.Fprintf(res, "event: redeem-fin\ndata: %s\n\n", metaJsonData)
	flusher.Flush()
	return nil
}

func redeem(fid, kid int, code string) (model.RedeemResponse, error) {
	payload := model.RedeemBodyKS{
		Fid: fid,
		Kid: kid,
		Cdk: code,
	}
	return postSigned(REDEEM_URL, payload)
}

func signPayload(data model.RedeemBodyKS) map[string]string {
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

func postSigned(url string, payload model.RedeemBodyKS) (model.RedeemResponse, error) {
	payload.Time = fmt.Sprintf("%d", time.Now().Unix())

	signed := signPayload(payload)

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

func chunkArray(arr []db.Player, size int) [][]db.Player {
	var result [][]db.Player
	for i := 0; i < len(arr); i += size {
		end := min(i+size, len(arr))
		result = append(result, arr[i:end])
	}
	return result
}

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Second)
}
