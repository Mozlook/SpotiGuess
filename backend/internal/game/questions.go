package game

import (
	"backend/internal/lastfm"
	"backend/internal/model"
	"backend/internal/spotify"
	"fmt"
	"log"
	"math/rand"
)

// GenerateQuestions generates quiz questions from a list of Spotify tracks.
//
// It expects:
//   - a slice of model.Track structs containing metadata about tracks,
//   - an OAuth access token to use for Spotify fallback search.
//
// The function performs the following steps for each track:
//
// 1. Skips the track if it has an empty ID.
//
// 2. Attempts to fetch 3 similar track titles using the Last.fm API.
//
//  3. If Last.fm fails (either by error or empty result), it falls back to
//     Spotify's search API using SimiliarFallback() to generate distractor answers.
//
// 4. If both methods fail to provide alternatives, the track is skipped.
//
//  5. Calculates a randomized playback start position for the track,
//     choosing a moment between 0 and (duration - 15 seconds), ensuring
//     there's at least a 15-second buffer from the end.
//
// 6. Constructs a model.Question object:
//   - Adds the correct track name along with 3 distractor titles
//   - Shuffles the answer options
//   - Assigns a unique ID ("q1", "q2", etc)
//   - Includes the playback position (in milliseconds)
//
// 7. Appends the question to the final result list.
//
// Returns:
//   - A slice of model.Question ready for the quiz,
//   - Or an error if critical steps fail.
//
// Example response:
//
//	[
//	  {
//	    "id": "q1",
//	    "trackId": "abc123",
//	    "trackName": "Shape of You",
//	    "positionMs": 90213,
//	    "answerOptions": ["Shape of You", "Thinking Out Loud", "Photograph", "Castle on the Hill"],
//	    "correctAnswer": "Shape of You"
//	  },
//	  ...
//	]
func GenerateQuestions(tracks []model.Track, token string) ([]model.Question, error) {
	var questions []model.Question
	for i, track := range tracks {
		var question model.Question

		if track.ID == "" {
			log.Printf("Skipping empty track ID for %s", track.Name)
			continue
		}

		recommendations, err := lastfm.FetchSimilar(track)
		if err != nil || len(recommendations) == 0 {
			log.Printf("Last.fm failed for track %s: %v — trying fallback", track.ID, err)
			recommendations, err = spotify.SimiliarFallback(track, token)
			if err != nil || len(recommendations) == 0 {
				log.Printf("Fallback also failed for track %s: %v", track.ID, err)
				continue
			}
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
