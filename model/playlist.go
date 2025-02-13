package model

import (
	"strconv"
	"time"

	"github.com/navidrome/navidrome/model/criteria"
	"github.com/navidrome/navidrome/utils"
)

type Playlist struct {
	ID        string         `structs:"id" json:"id"          orm:"column(id)"`
	Name      string         `structs:"name" json:"name"`
	Comment   string         `structs:"comment" json:"comment"`
	Duration  float32        `structs:"duration" json:"duration"`
	Size      int64          `structs:"size" json:"size"`
	SongCount int            `structs:"song_count" json:"songCount"`
	OwnerName string         `structs:"-" json:"ownerName"`
	OwnerID   string         `structs:"owner_id" json:"ownerId"  orm:"column(owner_id)"`
	Public    bool           `structs:"public" json:"public"`
	Tracks    PlaylistTracks `structs:"-" json:"tracks,omitempty"`
	Path      string         `structs:"path" json:"path"`
	Sync      bool           `structs:"sync" json:"sync"`
	CreatedAt time.Time      `structs:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `structs:"updated_at" json:"updatedAt"`

	// SmartPlaylist attributes
	Rules       *criteria.Criteria `structs:"-" json:"rules"`
	EvaluatedAt time.Time          `structs:"evaluated_at" json:"evaluatedAt"`
}

func (pls Playlist) IsSmartPlaylist() bool {
	return pls.Rules != nil && pls.Rules.Expression != nil
}

func (pls Playlist) MediaFiles() MediaFiles {
	if len(pls.Tracks) == 0 {
		return nil
	}
	return pls.Tracks.MediaFiles()
}

func (pls *Playlist) RemoveTracks(idxToRemove []int) {
	var newTracks PlaylistTracks
	for i, t := range pls.Tracks {
		if utils.IntInSlice(i, idxToRemove) {
			continue
		}
		newTracks = append(newTracks, t)
	}
	pls.Tracks = newTracks
}

func (pls *Playlist) AddTracks(mediaFileIds []string) {
	pos := len(pls.Tracks)
	for _, mfId := range mediaFileIds {
		pos++
		t := PlaylistTrack{
			ID:          strconv.Itoa(pos),
			MediaFileID: mfId,
			MediaFile:   MediaFile{ID: mfId},
			PlaylistID:  pls.ID,
		}
		pls.Tracks = append(pls.Tracks, t)
	}
}

func (pls *Playlist) AddMediaFiles(mfs MediaFiles) {
	pos := len(pls.Tracks)
	for _, mf := range mfs {
		pos++
		t := PlaylistTrack{
			ID:          strconv.Itoa(pos),
			MediaFileID: mf.ID,
			MediaFile:   mf,
			PlaylistID:  pls.ID,
		}
		pls.Tracks = append(pls.Tracks, t)
	}
}

type Playlists []Playlist

type PlaylistRepository interface {
	ResourceRepository
	CountAll(options ...QueryOptions) (int64, error)
	Exists(id string) (bool, error)
	Put(pls *Playlist) error
	Get(id string) (*Playlist, error)
	GetWithTracks(id string) (*Playlist, error)
	GetAll(options ...QueryOptions) (Playlists, error)
	FindByPath(path string) (*Playlist, error)
	RefreshStatus(playlistId string) error
	Delete(id string) error
	Tracks(playlistId string) PlaylistTrackRepository
}

type PlaylistTrack struct {
	ID          string `json:"id"          orm:"column(id)"`
	MediaFileID string `json:"mediaFileId" orm:"column(media_file_id)"`
	PlaylistID  string `json:"playlistId" orm:"column(playlist_id)"`
	MediaFile
}

type PlaylistTracks []PlaylistTrack

func (plt PlaylistTracks) MediaFiles() MediaFiles {
	mfs := make(MediaFiles, len(plt))
	for i, t := range plt {
		mfs[i] = t.MediaFile
	}
	return mfs
}

type PlaylistTrackRepository interface {
	ResourceRepository
	GetAll(options ...QueryOptions) (PlaylistTracks, error)
	Add(mediaFileIds []string) (int, error)
	AddAlbums(albumIds []string) (int, error)
	AddArtists(artistIds []string) (int, error)
	AddDiscs(discs []DiscID) (int, error)
	Delete(id ...string) error
	Reorder(pos int, newPos int) error
}
