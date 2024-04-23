package service

import (
	"TTMS_go/ttms/models"
	"errors"
	"fmt"
	"time"
)

func aviliable(movie models.Movie) error {
	if movie.Name == "" {
		return errors.New("movie name is empty")
	}
	if movie.Director == "" {
		return errors.New("movie director is empty")
	}
	if len(movie.Actor) == 0 {
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

func isTimeable(t *models.Theatre, play models.Play) error {
	fmt.Println("ssss")
	p := t.Plays
	q := p
	for p != nil {
		if p.Play.BeginTime.After(play.BeginTime) {
			break
		}
		q = p
		p = p.Next
	}
	fmt.Println("aaaa")
	if q != nil && q.Play.EndTime.Add(time.Minute*15).After(play.BeginTime) {
		return errors.New("演出开始时间早于上一场结束时间。")
	}
	if q != nil && q.Play.BeginTime.Before(play.EndTime.Add(time.Minute*15)) {
		return errors.New("演出结束时间与下一场开始时间冲突")
	}
	tmp := &models.Node{Play: play}
	fmt.Println("qqqq")
	if q != nil {
		q.Next = tmp
		tmp.Next = p
	}

	return nil
}
