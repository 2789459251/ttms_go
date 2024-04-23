package service

import (
	"TTMS_go/ttms/models"
	"encoding/json"
	"errors"
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

type playInfo struct {
	Begin time.Time
	End   time.Time
	Id    uint
}

func isTimeable(t *models.Theatre, play models.Play) error {
	cachePlay := t.Plays
	if len(cachePlay) == 0 {
		t.Plays = append(t.Plays, play)
		return nil
	}
	if len(cachePlay) == 1 {
		if play.BeginTime.After(cachePlay[0].EndTime.Add(15 * time.Minute)) {
			t.Plays = append(t.Plays, play)
			return nil
		}
		if play.EndTime.Add(15 * time.Minute).Before(cachePlay[0].BeginTime) {
			t.Plays = append([]models.Play{play}, cachePlay[0])
			return nil
		}

		jsonstr, _ := json.Marshal(playInfo{Begin: cachePlay[0].BeginTime, End: cachePlay[0].EndTime, Id: cachePlay[0].ID})
		return errors.New("添加的时间与原剧目安排冲突:" + string(jsonstr))
	}
	time1 := play.BeginTime
	time2 := play.EndTime
	play1, play2 := 0, 0
	for i := 1; i < len(cachePlay); i++ {
		if cachePlay[i-1].EndTime.Before(time1) && cachePlay[i].BeginTime.After(time2) {
			play1 = i - 1
			play2 = i
			break
		}
	}
	if cachePlay[play1].EndTime.Add(15*time.Minute).Before(time1) && time2.Add(15*time.Minute).Before(cachePlay[play2].BeginTime) {
		t.Plays = append(t.Plays[:play1+1], play)
		t.Plays = append(t.Plays, cachePlay[play2:]...)
		return nil
	}
	jsonstr1, _ := json.Marshal(playInfo{Begin: cachePlay[play1].BeginTime, End: cachePlay[play1].EndTime, Id: cachePlay[play1].ID})
	jsonstr2, _ := json.Marshal(playInfo{Begin: cachePlay[play2].BeginTime, End: cachePlay[play2].EndTime, Id: cachePlay[play2].ID})
	return errors.New("新增剧目与原剧目时间冲突:" + string(jsonstr1) + " " + string(jsonstr2))
}
