package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LikesClient struct {
	BaseURL string
}

type Like struct {
	LikeID     uint   `json:"like_id"`
	SourceID   uint   `json:"source_id"`
	TargetID   uint   `json:"target_id"`
	TargetType string `json:"target_type"`
}

type LikesResponse struct {
	Likes []Like `json:"likes"`
}

func (c *LikesClient) GetMyLikes(token string) ([]uint, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/likes/myLikes", nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("likes service returned status: %d", resp.StatusCode)
	}

	var likesResp LikesResponse
	if err := json.NewDecoder(resp.Body).Decode(&likesResp); err != nil {
		return nil, err
	}

	tripIDs := make([]uint, 0)
	for _, like := range likesResp.Likes {
		if like.TargetType == "trip" {
			tripIDs = append(tripIDs, like.TargetID)
		}
	}

	return tripIDs, nil
}