package service

import (
	"TTMS_go/ttms/models"
	"errors"
)

func aviliable(movie models.Movie) error {
	if movie.Name == "" {
		return errors.New("movie name is empty")
	}
	if movie.Director == "" {
		return errors.New("movie director is empty")
	}
	if movie.Actor == "" {
		return errors.New("movie actor is empty")
	}
	if movie.Duration == 0 {
		return errors.New("movie duration is empty")
	}
	if movie.ReleaseTime.IsZero() {
		return errors.New("movie release time is empty")
	}
	return nil
}
