package tibia

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var ErrNoMatches = fmt.Errorf("no matches")

type Photo struct {
	file     *os.File
	fileName string
	time     time.Time
	name     string
	_type    photoType
	number   int
}

func (p *Photo) FileName() string {
	return filepath.Base(p.fileName)
}

// FullPath rename
func (p *Photo) FullPath() string {
	return fmt.Sprintf("%s/%s/%s", p.name, p._type, p.FileName())
}

func (p *Photo) Album() string {
	return fmt.Sprintf("%s/%s", p.name, p._type)
}

func (p *Photo) Keep() bool {
	if p._type.Is(LevelUp, SkillUp, HighestDamageDealt) && p.number != 6 {
		return false
	}

	return true
}

func (p *Photo) File() io.Reader {
	return p.file
}

func (p *Photo) Time() time.Time {
	return p.time
}

var fileNameRegexp = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})_(\d{9})_(.*?)_(.*?)_?(\d)?\.png`)

func ParsePhotoFromFileName(name string) (*Photo, error) {
	p := &Photo{
		fileName: name,
	}
	matches := fileNameRegexp.FindStringSubmatch(name)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%w found in: %s", ErrNoMatches, name)
	}

	if len(matches) != 6 {
		return nil, fmt.Errorf("found not enough matches in: %s", name)
	}

	p.name = matches[3]
	p._type = decodePhotoType(matches[4])
	if !p._type.IsValid() {
		return nil, fmt.Errorf("photo type %s does not exist", p._type)
	}
	var err error
	p.time, err = parseTime(matches[1] + matches[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	p.number = 6

	if matches[5] != "" {
		n, err := strconv.Atoi(matches[5])
		if err != nil {
			return nil, fmt.Errorf("failed to convert string to number: %w", err)
		}
		p.number = n
	}
	return p, nil
}

func ParsePhotoFromFile(file *os.File) (*Photo, error) {
	p, err := ParsePhotoFromFileName(file.Name())
	if err != nil {
		return nil, err
	}
	p.file = file
	return p, nil
}

var loc *time.Location

func parseTime(s string) (time.Time, error) {
	s = s[:len(s)-3] + "." + s[:3] //fix missing dot for fraction of seconds
	if loc == nil {
		var err error
		loc, err = time.LoadLocation("Europe/Stockholm") //TODO: fix some env variable here
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to load loc: %w", err)
		}
	}
	return time.ParseInLocation("2006-01-02150405.999", s, loc)
}
