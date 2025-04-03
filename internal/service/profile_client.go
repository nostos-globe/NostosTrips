package service

import (
    "encoding/json"
    "fmt"
    "net/http"
	"main/internal/models"
)

type ProfileClient struct {
    BaseURL string
}


type FollowResponse struct {
    Follow struct {
        Count    int       `json:"count"`
        Profiles []models.Profile `json:"profiles"`
    } `json:"Follow"`
}

func (c *ProfileClient) GetFollowing(token string, userID uint) ([]int, error) {
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%d/following", c.BaseURL, userID), nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Cookie", "auth_token="+token)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get following users: %d", resp.StatusCode)
    }

    var response FollowResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }

    // Extract just the profile IDs
    var followingIDs []int
    for _, profile := range response.Follow.Profiles {
        followingIDs = append(followingIDs, int(profile.ProfileID))
    }

    return followingIDs, nil
}

func (c *ProfileClient) GetFollowers(token string, userID uint) ([]int, error) {
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%d/followers", c.BaseURL, userID), nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Cookie", "auth_token="+token)
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to get followers: %d", resp.StatusCode)
    }

    var response FollowResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }

    // Extract just the profile IDs
    var followerIDs []int
    for _, profile := range response.Follow.Profiles {
        followerIDs = append(followerIDs, int(profile.ProfileID))
    }

    return followerIDs, nil
}