package player

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
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
func AddPlayer(c *echo.Context) error {
	KSFetchAvailable := global.KSFetchAvailable
	dbApp := db.InitDB()
	defer dbApp.Pool.Close()

	ctx := context.Background()

	// pid from params
	body := new(model.PlayerReqBody)
	// Step 2: Bind JSON body into struct
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
	}

	// validating body
	v := validator.New()
	if err := v.Struct(body); err != nil {
		// Instead of dumping raw error, let’s format field-specific messages
		errs := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			field := e.Field()
			tag := e.Tag()

			switch field {
			case "Pid":
				errs[field] = "Must be a valid PlayerID"
			case "Kid":
				errs[field] = "Must give a valid KingdomID (>=1)"
			default:
				// fallback based on tag
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
		return c.JSON(http.StatusBadRequest, errs)
	}

	// parse pid and alliance
	pidInt, err := strconv.Atoi(body.Pid)
	if err != nil {
		slog.Error("Invalid pid", "pid", body.Pid, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid pid"})
	}

	var allianceParam pgtype.Text
	if body.Alliance == "" {
		// leave it invalid so Postgres default kicks in
		allianceParam = pgtype.Text{String: "DES", Valid: true}
	} else {
		allianceParam = pgtype.Text{String: body.Alliance, Valid: true}
	}

	// Grabing the player data from a API (Not Available rn)
	if KSFetchAvailable {
		url := "https://kingshot.net/api/player-info"
		client := resty.New()
		defer client.Close()
		var pD model.KSFetchPlayerResponse
		_, e := client.R().SetResult(pD).SetRetryCount(3).SetHeader("Content-Type", "application/json").SetQueryParam("playerId", body.Pid).Post(url)

		if e != nil {
			slog.Error("Error in fetching Player", "err", e)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unable to fetch player"})
		}

		p, e := dbApp.Queries.PushPlayer(ctx, db.PushPlayerParams{
			Pid:   int32(pidInt),
			Kid:   int32(pD.Data.Kingdom),
			Dname: pgtype.Text{String: pD.Data.Name, Valid: true},
			Pfp: pgtype.Text{
				String: pD.Data.ProfilePhoto,
				Valid:  true,
			},
			Alliance: allianceParam,
		})

		if e != nil {
			slog.Error("Error in AddPlayer Route", "err", e)
			c.Response().Header().Set("Content-Type", "application/json")
			c.Response().Header().Set("Cache-Control", "no-store")
			return c.JSON(500, map[string]string{"message": "Unable to add the player"})
		}
		return c.JSON(http.StatusCreated, model.PlayerResponse{Message: "Created Player", Type: 1, Player: model.PlayerInfo{
			Pid:      p.Pid,
			Dname:    p.Dname.String,
			Kid:      p.Kid,
			Alliance: p.Alliance.String,
			Pfp:      p.Pfp.String,
		}})
	}

	if err != nil {
		slog.Error("Invalid kid", "kid", body.Kid, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid kid"})
	}

	p, e := dbApp.Queries.PushPlayer(ctx, db.PushPlayerParams{
		Pid:      int32(pidInt),
		Kid:      int32(body.Kid),
		Alliance: allianceParam,
	})
	if e != nil {
		slog.Error("Error in AddPlayer Route", "err", e)
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().Header().Set("Cache-Control", "no-store")
		return c.JSON(500, map[string]string{"message": "Unable to add the player"})
	}

	return c.JSON(http.StatusCreated, model.PlayerResponse{Message: "Created Player without API", Type: 2, Player: model.PlayerInfo{
		Pid:      p.Pid,
		Dname:    p.Dname.String,
		Kid:      p.Kid,
		Alliance: p.Alliance.String,
		Pfp:      p.Pfp.String,
	}})
}
