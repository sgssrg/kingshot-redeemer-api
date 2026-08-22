package redeem

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"resty.dev/v3"
)

func RedeemSingle(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {

		client := resty.New()
		defer client.Close()

		body, err := bindAndValidatePlayer(c)
		if err != nil {
			return err
		}

		pidInt, err := strconv.Atoi(body.Pid)
		if err != nil {
			slog.Error("Invalid pid", "pid", body.Pid, "err", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid pid"})
		}
		var pD model.PlayerInfo
		url := os.Getenv("host") + "/scraper/player/" + body.Pid
		_, err = client.R().SetResult(&pD).SetRetryCount(3).SetHeader("Content-Type", "application/json").Get(url)

		resp, err := Redeem(int(pD.Pid), int(pD.Kid), body.Giftcode)
		redeemEvent := model.RedeemSSE{ // Redeem Response
			Fid:       strconv.Itoa(int(pidInt)),
			Code:      body.Giftcode,
			RawResult: resp,
			Success:   err == nil,
			Time:      time.Now().Format(time.RFC3339),
		}
		slog.Info("resp", "data", resp)
		if err == nil && resp.ErrCode == 0 && resp.Msg == "" || resp.Code == 40004 && resp.Msg == "TIMEOUT RETRY." {
			slog.Warn("Received empty response, retrying in 5 seconds...", "fid", int(pD.Pid))
			sleep(5)
			resp, err = Redeem(int(pD.Pid), 1420, body.Giftcode)
		}
		switch redeemEvent.RawResult.ErrCode {
		case 20000:
			redeemEvent.RefinedResult.Message = "Redeemed By Bot 💝"
			redeemEvent.RefinedResult.Redeemed = true
			redeemEvent.RefinedResult.RedeemedByBot = true
			redeemEvent.RefinedResult.TypeRedeemedCode = 0
		case 40004:
			redeemEvent.RefinedResult.Message = "Timeout Retry 🧱"
			redeemEvent.RefinedResult.Redeemed = false
			redeemEvent.RefinedResult.RedeemedByBot = false
			redeemEvent.RefinedResult.TypeRedeemedCode = 1
		case 40008:
			redeemEvent.RefinedResult.Message = "Manual ❤️‍🩹"
			redeemEvent.RefinedResult.Redeemed = true
			redeemEvent.RefinedResult.RedeemedByBot = false
			redeemEvent.RefinedResult.TypeRedeemedCode = 2
		case 40020:
			redeemEvent.RefinedResult.Message = "Dead PlayerID 🛠️"
			redeemEvent.RefinedResult.Redeemed = false
			redeemEvent.RefinedResult.RedeemedByBot = false
			redeemEvent.RefinedResult.TypeRedeemedCode = 3
		case 40005, 40007, 40014:
			redeemEvent.RefinedResult.Message = "Expired/Dead GiftCode 🗑️"
			redeemEvent.RefinedResult.Redeemed = false
			redeemEvent.RefinedResult.RedeemedByBot = false
			redeemEvent.RefinedResult.TypeRedeemedCode = 4
		}
		c.JSON(http.StatusAccepted, redeemEvent)
		return nil
	}
}

func bindAndValidatePlayer(c *echo.Context) (*model.GiftCodeReqBody, error) {
	body := new(model.GiftCodeReqBody)
	if err := c.Bind(body); err != nil {
		return nil, c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		return nil, c.JSON(http.StatusBadRequest, buildValidationErrors(err.(validator.ValidationErrors)))
	}

	return body, nil
}

func buildValidationErrors(validationErrs validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, e := range validationErrs {
		field := e.Field()
		tag := e.Tag()

		switch field {
		case "Pid":
			errs[field] = "Must be a valid PlayerID"
		case "Giftcode":
			errs[field] = "Must provide a Giftcode to Proceed"
		default:
			switch tag {
			case "required":
				errs[field] = "This field is required"
			case "gte":
				errs[field] = "Value too small"
			case "lte":
				errs[field] = "Value too large"
			default:
				errs[field] = "Validation failed on " + tag
			}
		}
	}
	return errs
}
