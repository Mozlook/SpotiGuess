package game

import (
	"backend/internal/lastfm"
	"backend/internal/model"
	"fmt"
	"log"
	"math/rand"
)

// GenerateQuestions generates a slice of quiz questions from the given track list.
//
// Each question includes a correct answer (the original track) and 3 fake answers
// fetched using the Last.fm API via the FetchSimilar function.
//
// It performs the following steps:
//
//  1. Iterates over the provided slice of model.Track.
//  2. For each track:
//     - Skips it if the track has an empty ID (to avoid invalid entries).
//     - Calls lastfm.FetchSimilar to obtain up to 3 similar track names.
//     - Constructs a model.Question:
//     • ID: "q1", "q2", ...
//     • TrackID: original track's ID
//     • TrackName: original track name
//     • AnswerOptions: 3 fake answers + correct answer, in random order
//     • CorrectAnswer: original track name
//  3. Appends the generated question to the result slice.
//
// Returns:
//   - []model.Question on success
//   - error if any call to FetchSimilar fails
//
// Logging is used to report skipped tracks or failures to fetch recommendations.
//
// Example question structure:
//
//	{
//	  "id": "q1",
//	  "trackId": "abc123",
//	  "trackName": "Blinding Lights",
//	  "answerOptions": ["Starboy", "Can't Feel My Face", "Save Your Tears", "Blinding Lights"],
//	  "correctAnswer": "Blinding Lights"
//	}
func GenerateQuestions(tracks []model.Track) ([]model.Question, error) {
	var questions []model.Question
	for i, track := range tracks {
		var question model.Question

		if track.ID == "" {
			log.Printf("Skipping empty track ID for %s", track.Name)
			continue
		}

		recommendations, err := lastfm.FetchSimilar(track)
		if err != nil {
			log.Printf("Failed to fetch recommendations for track %s: %v", track.ID, err)
			return nil, err
		}

		trackDuration := track.Duration // w ms
		maxStart := trackDuration - 15000
		maxStart = max(maxStart, 0)
		startMs := rand.Intn(maxStart)

		question.ID = fmt.Sprintf("q%d", i+1)
		question.TrackID = track.ID
		question.TrackName = track.Name
		question.AnswerOptions = recommendations
		question.AnswerOptions = append(question.AnswerOptions, track.Name)
		question.CorrectAnswer = track.Name
		question.PositionMs = startMs
		rand.Shuffle(len(question.AnswerOptions), func(i, j int) {
			question.AnswerOptions[i], question.AnswerOptions[j] = question.AnswerOptions[j], question.AnswerOptions[i]
		})

		questions = append(questions, question)
	}

	return questions, nil

}
