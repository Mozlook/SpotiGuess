package game

import (
	"backend/internal/spotify"
	"fmt"
	"math/rand/v2"
)

type Question struct {
	ID            string   `json:"id"`
	TrackID       string   `json:"trackId"`
	TrackName     string   `json:"trackName"`
	AnswerOptions []string `json:"options"`
	CorrectAnswer string   `json:"correct"`
}
type Track struct {
	ID      string
	Name    string
	Artists []string
}

func GenerateQuestions(tracks []Track, token string) ([]Question, error) {
	var questions []Question
	for i, track := range tracks {
		var question Question
		recommendations, err := spotify.FetchRecommendations(track.ID, token)
		if err != nil {
			return nil, err
		}
		question.ID = fmt.Sprintf("q%d", i+1)
		question.TrackID = track.ID
		question.TrackName = track.Name
		question.AnswerOptions = recommendations
		question.AnswerOptions = append(question.AnswerOptions, track.Name)
		question.CorrectAnswer = track.Name
		rand.Shuffle(len(question.AnswerOptions), func(i, j int) {
			question.AnswerOptions[i], question.AnswerOptions[j] = question.AnswerOptions[j], question.AnswerOptions[i]
		})

		questions = append(questions, question)
	}

	return questions, nil
}
