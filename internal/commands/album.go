package commands

import (
	"context"
	"fmt"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/query"

	"github.com/urfave/cli"
)

var AlbumCommand = cli.Command{
	Name:    "album",
	Aliases: []string{""},
	Usage:   "Adds photos to an album",
	Flags:   albumFlags,
	Action:  albumAction,
}

var albumFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "album",
		Value: "",
		Usage: "Album to add photos to",
	},
}

// albumAction adds photos to an album.
func albumAction(ctx *cli.Context) error {
	conf := config.NewConfig(ctx)

	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := conf.Init(cctx); err != nil {
		return err
	}

	conf.MigrateDb()

	db := conf.Db()
	q := query.New(db)
	albumResults, err := q.Albums(form.AlbumSearch{Name: ctx.String("album")})

	if err != nil {
		log.Errorf("album: %s", err)
		return nil
	}

	var albumName string
	var albumUUID string
	if len(albumResults) > 0 {
		a := albumResults[0]

		albumName = a.AlbumName
		albumUUID = a.AlbumUUID
	} else {
		// try creating the album if not found
		a := entity.NewAlbum(ctx.String("album"))

		if err := db.Create(a).Error; err != nil {
			log.Errorf("create album: %s", err)
			return nil
		}

		albumName = a.AlbumName
		albumUUID = a.AlbumUUID
		log.Infof("new album created %s\n", albumName)
	}

	photos, _, err := q.PhotosByFilenames(ctx.Args())

	var added []*entity.PhotoAlbum

	for _, p := range photos {
		added = append(added, entity.NewPhotoAlbum(p.PhotoUUID, albumUUID).FirstOrCreate(db))
	}

	if len(added) == 1 {
		log.Info(fmt.Sprintf("one photo added to %s", albumName))
	} else {
		log.Info(fmt.Sprintf("%d photos added to %s", len(added), albumName))
	}
	return nil
}
