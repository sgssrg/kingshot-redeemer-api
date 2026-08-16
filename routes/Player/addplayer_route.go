package player

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"resty.dev/v3"

	"gitlab.com/ribonin/apis/kingshot-redeem/db"
	"gitlab.com/ribonin/apis/kingshot-redeem/global"
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
)

// AddPlayer registers a new player.
// @Summary Add a player
// @Description Adds a new player to the database. When KSFetchAvailable is enabled, profile data is fetched from the KingShot API; otherwise only the provided fields are stored.
// @Tags Player
// @Accept json
// @Produce json
// @Param player body model.PlayerReqBody true "Player to add"
// @Success 201 {object} model.PlayerAPIResponse "Created Player examples" examples(
//
//	{"Created Player without Player Fetch": {"message":"Created Player without API","type":2,"player":{"pid":123,"kid":789,"alliance":"XYZ"}}},
//	{"Created Player using Player Fetch": {"message":"Created Player","type":1,"player":{"pid":201927560,"kid":1420,"dName":"name-of-the-profile","pfp":"url-of-pfp","alliance":"XYZ"}}}
//
// )
// @Failure 400 {object} map[string]interface{} "Invalid JSON or validation error"
// @Failure 500 {object} map[string]interface{} "Unable to add the player"
// @Router /player/add [post]
func AddPlayer(dbApp *db.App) echo.HandlerFunc {
	return func(c *echo.Context) error {
		KSFetchAvailable := global.KSFetchAvailable
		StratForgeFetchAvailable := global.StratForgeFetchAvailable

		ctx := context.Background()
		client := resty.New()
		defer client.Close()

		body, err := bindAndValidatePlayer(c)
		if err != nil {
			return err
		}

		pidInt, err := strconv.Atoi(body.Pid)
		if err != nil {
			slog.Error("Invalid pid", "pid", body.Pid, "err", err)
			return c.JSON(http.StatusBadRequest, model.AddPlayerValidator{WrongField: "pid", Message: "Invalid PiD & It should be Interger > 0"})
		}

		allianceParam := buildAllianceParam(body.Alliance)

		if KSFetchAvailable {
			return handleKSFetch(ctx, c, client, dbApp, body, pidInt, allianceParam)
		}

		if StratForgeFetchAvailable {
			handled, handlerErr := tryAddWithStratForge(ctx, c, client, dbApp, body, pidInt, allianceParam)
			if handled || handlerErr != nil {
				return handlerErr
			}
		}

		return addPlayerManual(ctx, c, dbApp, body, pidInt, allianceParam)
	}
}

func bindAndValidatePlayer(c *echo.Context) (*model.PlayerReqBody, error) {
	body := new(model.PlayerReqBody)
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
		case "Kid":
			errs[field] = "Must give a valid KingdomID (>=1)"
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

func buildAllianceParam(alliance string) sql.NullString {
	if alliance == "" {
		return sql.NullString{String: "DES", Valid: true}
	}
	return sql.NullString{String: alliance, Valid: true}
}

func handleKSFetch(ctx context.Context, c *echo.Context, client *resty.Client, dbApp *db.App, body *model.PlayerReqBody, pidInt int, allianceParam sql.NullString) error {
	var pD model.KSFetchPlayerResponse
	_, err := client.R().SetResult(&pD).SetRetryCount(3).SetHeader("Content-Type", "application/json").SetQueryParam("playerId", body.Pid).Post("https://kingshot.net/api/player-info")
	if err != nil {
		slog.Error("Error in fetching Player", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to fetch player"})
	}

	p, err := dbApp.Queries.PushPlayer(ctx, db.PushPlayerParams{
		Pid:      int64(pidInt),
		Kid:      int64(pD.Data.Kingdom),
		Dname:    sql.NullString{String: pD.Data.Name, Valid: true},
		Pfp:      sql.NullString{String: pD.Data.ProfilePhoto, Valid: true},
		Alliance: allianceParam,
	})
	if err != nil {
		slog.Error("Error in AddPlayer Route", "err", err)
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().Header().Set("Cache-Control", "no-store")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to add the player"})
	}

	return c.JSON(http.StatusCreated, model.PlayerResponse{Message: "Created Player", Type: 1, Player: model.PlayerInfo{
		Pid:      uint(p.Pid),
		Dname:    p.Dname.String,
		Kid:      uint(p.Kid),
		Alliance: p.Alliance.String,
		Pfp:      p.Pfp.String,
	}})
}

func tryAddWithStratForge(ctx context.Context, c *echo.Context, client *resty.Client, dbApp *db.App, body *model.PlayerReqBody, pidInt int, allianceParam sql.NullString) (bool, error) {
	var pD model.PlayerInfo
	url := os.Getenv("host") + "/scraper/player/" + body.Pid
	pF, err := client.R().SetResult(&pD).SetRetryCount(3).SetHeader("Content-Type", "application/json").Get(url)
	if err != nil {
		slog.Error("Fetch failed", "error", err)
		return true, c.JSON(http.StatusInternalServerError, "fetch error")
	}

	if pF.StatusCode() != 202 {
		return false, nil
	}

	p, err := dbApp.Queries.PushPlayer(ctx, db.PushPlayerParams{
		Pid:      int64(pidInt),
		Kid:      int64(pD.Kid),
		Dname:    sql.NullString{String: pD.Dname, Valid: true},
		Pfp:      sql.NullString{String: pD.Pfp, Valid: true},
		Alliance: sql.NullString{String: pD.Alliance, Valid: true},
	})
	if err != nil {
		slog.Error("Error in AddPlayer Route", "err", err)
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().Header().Set("Cache-Control", "no-store")
		return true, c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to add the player", "error": err.Error()})
	}

	slog.Info("Player Added with PiD - " + body.Pid)
	return true, c.JSON(http.StatusCreated, model.PlayerResponse{Message: "Created Player with self API", Type: 2, Player: model.PlayerInfo{
		Pid:      uint(p.Pid),
		Dname:    p.Dname.String,
		Kid:      uint(p.Kid),
		Alliance: p.Alliance.String,
		Pfp:      p.Pfp.String,
	}})
}

func addPlayerManual(ctx context.Context, c *echo.Context, dbApp *db.App, body *model.PlayerReqBody, pidInt int, allianceParam sql.NullString) error {
	if body.Kid == nil {
		return c.JSON(http.StatusBadRequest, model.AddPlayerValidator{
			WrongField: "kid",
			Message:    "Kid Should be provided as both StratForge Scraper and KSNet isn't available",
		})
	}

	p, err := dbApp.Queries.PushPlayer(ctx, db.PushPlayerParams{
		Pid:      int64(pidInt),
		Kid:      int64(*body.Kid),
		Alliance: allianceParam,
	})
	if err != nil {
		slog.Error("Error in AddPlayer Route", "err", err)
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().Header().Set("Cache-Control", "no-store")
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to add the player"})
	}

	return c.JSON(http.StatusCreated, model.PlayerResponse{Message: "Created Player without API", Type: 3, Player: model.PlayerInfo{
		Pid:      uint(p.Pid),
		Dname:    p.Dname.String,
		Kid:      uint(p.Kid),
		Alliance: p.Alliance.String,
		Pfp:      p.Pfp.String,
	}})
}
