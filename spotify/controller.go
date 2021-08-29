package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PlaylistTrackResponse struct {
	Href  string `json:"href"`
	Items []struct {
		AddedAt time.Time `json:"added_at"`
		AddedBy struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"added_by"`
		IsLocal      bool        `json:"is_local"`
		PrimaryColor interface{} `json:"primary_color"`
		Track        struct {
			Album struct {
				AlbumType string `json:"album_type"`
				Artists   []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					Height int    `json:"height"`
					URL    string `json:"url"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				TotalTracks          int    `json:"total_tracks"`
				Type                 string `json:"type"`
				URI                  string `json:"uri"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Episode          bool     `json:"episode"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href        string `json:"href"`
			ID          string `json:"id"`
			IsLocal     bool   `json:"is_local"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			Track       bool   `json:"track"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
		} `json:"track"`
		VideoThumbnail struct {
			URL interface{} `json:"url"`
		} `json:"video_thumbnail"`
	} `json:"items"`
	Limit    int         `json:"limit"`
	Next     interface{} `json:"next"`
	Offset   int         `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int         `json:"total"`
}

type Music struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Query  string `json:"query"`
	Image  string `json:"imgage"`
}

func GetPlaylistTracks(c *fiber.Ctx) error {
	var musics []Music
	var spotifyResp PlaylistTrackResponse
	accessToken := c.Get("access-token")
	playlistURI := c.Query("playlist_uri", "https://open.spotify.com/playlist/4nlyLwDykN13rB1IDHqNDl?si=cJCe_cAtR1WnuCVc3Lz3kg&dl_branch=1")
	uri, _ := url.Parse(playlistURI)
	playlistURI = path.Base(uri.Path)

	if accessToken == "" || len(accessToken) < 1 {
		return c.Status(http.StatusUnauthorized).JSON(struct {
			Code    int
			Message string
		}{
			Code:    http.StatusUnauthorized,
			Message: "missing access token from header",
		})
	}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistURI), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&spotifyResp)
	if err != nil {
		logrus.Error("error decoding payload: %v, err")
	}

	for _, item := range spotifyResp.Items {
		title := item.Track.Name
		artist := item.Track.Artists[0].Name

		musics = append(musics, Music{
			Title:  title,
			Artist: artist,
			Query:  fmt.Sprintf("%s %s", artist, title),
			Image:  item.Track.Album.Images[0].URL,
		})
	}
	return c.Status(http.StatusOK).JSON(struct {
		Code   int     `json:"code"`
		Tracks []Music `json:"tracks"`
	}{
		Code:   http.StatusOK,
		Tracks: musics,
	})
}
