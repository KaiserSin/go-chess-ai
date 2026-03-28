package assets

import (
	"bytes"
	"embed"
	"image/png"
	"io/fs"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed pieces/*.png
var pieceFiles embed.FS

type PieceSprites struct {
	sprites map[string]*ebiten.Image
}

func LoadPieceSprites() PieceSprites {
	sprites := PieceSprites{
		sprites: make(map[string]*ebiten.Image),
	}

	paths, err := fs.Glob(pieceFiles, "pieces/*.png")
	if err != nil {
		return sprites
	}

	for _, filePath := range paths {
		data, err := pieceFiles.ReadFile(filePath)
		if err != nil {
			continue
		}

		source, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			continue
		}

		key := strings.TrimSuffix(path.Base(filePath), ".png")
		sprites.sprites[key] = ebiten.NewImageFromImage(source)
	}

	return sprites
}

func (p PieceSprites) Lookup(assetKey string) (*ebiten.Image, bool) {
	sprite, ok := p.sprites[assetKey]
	return sprite, ok
}

func (p PieceSprites) Count() int {
	return len(p.sprites)
}
