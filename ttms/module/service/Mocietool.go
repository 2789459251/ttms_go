package service

import (
	"TTMS_go/ttms/models"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	playId := strings.Split(t.Plays, " ")
	fmt.Println("pp:", t.Plays)
	fmt.Println("dp:", playId)
	cachePlay := models.FindPlayByIds(playId)
	fmt.Println("len:", len(cachePlay))
	defer func() {
		fmt.Println(t.Plays)
	}()
	if len(cachePlay) == 0 {
		t.Plays = strconv.Itoa(int(play.ID))
		fmt.Println("00000", t.Plays)
		return nil
	}
	if len(playId) == 1 {
		if play.BeginTime.After(cachePlay[0].EndTime.Add(15 * time.Minute)) {
			t.Plays = playId[0] + " " + strconv.Itoa(int(play.ID))
			return nil
		}
		if play.EndTime.Add(15 * time.Minute).Before(cachePlay[0].BeginTime) {
			tplay := t.Plays
			t.Plays = strconv.Itoa(int(play.ID)) + " " + tplay
			return nil
		}

		jsonstr, _ := json.Marshal(playInfo{Begin: cachePlay[0].BeginTime, End: cachePlay[0].EndTime, Id: cachePlay[0].ID})
		return errors.New("添加的时间与原剧目安排冲突:" + string(jsonstr))
	}
	time1 := play.BeginTime
	time2 := play.EndTime
	play1, play2 := 0, 0
	for i := 1; i < len(playId); i++ {
		if cachePlay[i-1].EndTime.Before(time1) && cachePlay[i].BeginTime.After(time2) {
			play1 = i - 1
			play2 = i
			break
		}
	}
	if cachePlay[play1].EndTime.Add(15*time.Minute).Before(time1) && time2.Add(15*time.Minute).Before(cachePlay[play2].BeginTime) {
		a := append(playId[:play1+1], strconv.Itoa(int(play.ID)))
		b := append(a, playId[play2:]...)
		t.Plays = strings.Join(b, " ")
		return nil
	}
	jsonstr1, _ := json.Marshal(playInfo{Begin: cachePlay[play1].BeginTime, End: cachePlay[play1].EndTime, Id: cachePlay[play1].ID})
	jsonstr2, _ := json.Marshal(playInfo{Begin: cachePlay[play2].BeginTime, End: cachePlay[play2].EndTime, Id: cachePlay[play2].ID})
	return errors.New("新增剧目与原剧目时间冲突:" + string(jsonstr1) + " " + string(jsonstr2))
}
