package redeem

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
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
// @Param code query string false "Gift code" example(KS0803)
// @Success 200 {string} string "redeem-res and redeem-fin event" example(event: redeem-res\ndata: {"fid":"123","code":"KS0803","result":{"errCode":20000,"msg":"Redeemed"},"success":true,"time":"2026-08-07T22:05:40+05:30"}\n\n event: redeem-fin\ndata: {"redeemed":10,"manual":2,"playerDead":1,"codeExpired":3})
// @Router /redeem/db [post]
func RedeemWithDB(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {
		reqContext := c.Request().Context()
		res := c.Response()
		ctx := context.Background()
		players, err := dbApp.Queries.GetAllPlayers(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
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
			return c.JSON(http.StatusBadRequest, model.UpdatePlayerSSEValidator{WrongField: "code", Message: "Please enter a valid code."})
		}

		chunks := chunkArray(players, 5)
		var playerMetaRedeemResponse model.PlayerMetaRedeemResponse
		var finalRes []model.RedeemSSE
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-reqContext.Done():
					return
				case <-ticker.C:
					fmt.Fprintf(res, ": keep-alive\n\n")
					flusher.Flush()
				}
			}
		}()

		for _, batch := range chunks {
			for _, p := range batch {
				// Check if the client disconnected before doing work
				select {
				case <-reqContext.Done():
					slog.Info("Client disconnected, stopping process")
					return nil
				default:
				}

				resp, err := Redeem(int(p.Pid), int(p.Kid), giftCode)
				redeemEvent := model.RedeemSSE{ // Redeem Response
					Fid:       strconv.Itoa(int(p.Pid)),
					Code:      giftCode,
					RawResult: resp,
					Success:   err == nil,
					Time:      time.Now().Format(time.RFC3339),
				}
				slog.Info("resp", "data", resp)
				finalRes = append(finalRes, redeemEvent)
				if err == nil && resp.ErrCode == 0 && resp.Msg == "" || resp.Code == 40004 {
					slog.Warn("Received empty response or TIMEOUT_RETRY, retrying in 5 seconds...", "fid", int(p.Pid))
					sleep(5)
					resp, err = Redeem(int(p.Pid), 1420, giftCode)
				}
				switch redeemEvent.RawResult.ErrCode {
				case 20000:
					redeemEvent.RefinedResult.Message = "Redeemed By Bot 💝"
					redeemEvent.RefinedResult.Redeemed = true
					redeemEvent.RefinedResult.TypeRedeemedCode = 0
					playerMetaRedeemResponse.Redeemed++
				case 40008:
					redeemEvent.RefinedResult.Message = "Manual ❤️‍🩹"
					redeemEvent.RefinedResult.Redeemed = true
					redeemEvent.RefinedResult.TypeRedeemedCode = 1
					playerMetaRedeemResponse.Manual++
				case 40020:
					redeemEvent.RefinedResult.Message = "Dead PlayerID 🛠️"
					redeemEvent.RefinedResult.Redeemed = false
					redeemEvent.RefinedResult.TypeRedeemedCode = 3
					playerMetaRedeemResponse.PlayerDead++
				case 40005, 40007, 40014:
					redeemEvent.RefinedResult.Message = "Expired/Dead GiftCode 🗑️"
					redeemEvent.RefinedResult.Redeemed = false
					redeemEvent.RefinedResult.TypeRedeemedCode = 3
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
				sleep(1)
			}
			// Delay between batches
			sleep(2)
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
