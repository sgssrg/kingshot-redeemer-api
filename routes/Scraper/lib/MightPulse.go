package lib

import (
	"gitlab.com/ribonin/apis/kingshot-redeem/model"
	"resty.dev/v3"
)

func MPAlliancePlayerScaper(KSKid *string, KSAlliance *string) (*model.MightPulseScrapeAlliance, int, error) {
	client := resty.New()

	mpURL := "https://mightpulse.com/api/alliances/lookup?kid=" + *KSKid + "&slug=" + *KSAlliance
	var aD model.MightPulseScrapeAlliance    
	aF, err := client.R().SetResult(aD).SetRetryCount(3).SetHeader("Content-Type", "application/json").Get(mpURL)

	return &aD, aF.StatusCode(), err
}
