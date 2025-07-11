package game

import (
	"backend/internal/model"
	"backend/internal/spotify"
	"fmt"
	"math/rand/v2"
)

// GenerateQuestions builds a list of quiz questions based on Spotify track data.
//
// It accepts a list of tracks (each with ID and name) and a Spotify access token
// used to fetch recommendations for each track.
//
// For each input track, the function:
//
//  1. Calls FetchRecommendations to get 3 similar track names (used as incorrect answers).
//  2. Appends the correct track name to the answer options.
//  3. Shuffles the resulting list of 4 answers.
//  4. Creates a Question struct with:
//     - A unique ID ("q1", "q2", ...),
//     - Track ID (used later for playback),
//     - Track name (used for result display),
//     - A shuffled list of answer options,
//     - The correct answer as a string.
//
// The resulting list of questions can be used to power a full quiz round.
//
// Returns:
// - A slice of Question structs
// - Or an error if fetching recommendations fails for any track
func GenerateQuestions(tracks []spotify.Track, token string) ([]model.Question, error) {
	var questions []model.Question
	for i, track := range tracks {
		var question model.Question
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
