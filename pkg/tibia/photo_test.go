package tibia

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	p1 := "2020-08-28144154390"
	p2 := "2022-07-17093153856"
	fmt.Println(parseTime(p1))
	fmt.Println(time.Parse("2006-01-02150405.000", p1))
	fmt.Println(time.Parse("2006-01-02150405.999", p2))
}

func TestParsePhotoFromFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		verify   func(t *testing.T, photo *Photo)
	}{
		{
			name:     "LevelUp 6 image",
			fileName: "/some/random/place/2024-03-22_220250767_Xon Vanetta_LevelUp_6.png",
			verify: func(t *testing.T, photo *Photo) {
				assert.Equal(t, 6, photo.number)
				assert.Equal(t, LevelUp, photo._type)
				assert.Equal(t, "2024-03-22_220250767_Xon Vanetta_LevelUp_6.png", photo.FileName())
				assert.Equal(t, "Xon Vanetta/LevelUp/2024-03-22_220250767_Xon Vanetta_LevelUp_6.png", photo.FullPath())
				assert.Equal(t, "Xon Vanetta/LevelUp", photo.Album())
				assert.Equal(t, true, photo.Keep())
				assert.Equal(t, "2024-03-22 22:02:50.202 +0100 CET", photo.Time().String()) //TODO this might break since TZ
			},
		},
		{
			name:     "SkillUp 5 image",
			fileName: "/some/random/place/2024-03-22_220250767_Xon Vanetta_SkillUp_5.png",
			verify: func(t *testing.T, photo *Photo) {
				assert.Equal(t, 5, photo.number)
				assert.Equal(t, SkillUp, photo._type)
				assert.Equal(t, "2024-03-22_220250767_Xon Vanetta_SkillUp_5.png", photo.FileName())
				assert.Equal(t, "Xon Vanetta/SkillUp/2024-03-22_220250767_Xon Vanetta_SkillUp_5.png", photo.FullPath())
				assert.Equal(t, "Xon Vanetta/SkillUp", photo.Album())
				assert.Equal(t, false, photo.Keep())
				assert.Equal(t, "2024-03-22 22:02:50.202 +0100 CET", photo.Time().String()) //TODO this might break since TZ
			},
		},
		{
			name:     "DeathPvE 2 image",
			fileName: "/some/random/place/2024-03-22_220250767_Xon Vanetta_DeathPvE_2.png",
			verify: func(t *testing.T, photo *Photo) {
				assert.Equal(t, 2, photo.number)
				assert.Equal(t, DeathPvE, photo._type)
				assert.Equal(t, "2024-03-22_220250767_Xon Vanetta_DeathPvE_2.png", photo.FileName())
				assert.Equal(t, "Xon Vanetta/DeathPvE/2024-03-22_220250767_Xon Vanetta_DeathPvE_2.png", photo.FullPath())
				assert.Equal(t, "Xon Vanetta/DeathPvE", photo.Album())
				assert.Equal(t, true, photo.Keep())
				assert.Equal(t, "2024-03-22 22:02:50.202 +0100 CET", photo.Time().String()) //TODO this might break since TZ
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			photo, err := ParsePhotoFromFileName(test.fileName)
			assert.NoError(t, err)
			test.verify(t, photo)
		})
	}
}
