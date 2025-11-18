package search

import (
	"picstagsbot/internal/domain/model"

	tele "gopkg.in/telebot.v4"
)

func (h *SearchHandler) sendPhotosAsAlbums(c tele.Context, userID int64, photos []*model.Photo) {
	for i := 0; i < len(photos); i += albumSize {
		end := i + albumSize
		if end > len(photos) {
			end = len(photos)
		}

		batch := photos[i:end]
		var album tele.Album

		for _, p := range batch {
			album = append(album, &tele.Photo{
				File:    tele.File{FileID: p.TelegramID},
				Caption: p.Description,
			})
		}

		if err := c.SendAlbum(album); err != nil {
			for _, p := range batch {
				_ = c.Send(&tele.Photo{
					File:    tele.File{FileID: p.TelegramID},
					Caption: p.Description,
				})
			}
		}
	}
}
