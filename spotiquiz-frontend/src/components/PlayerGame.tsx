import type { Question } from "./GamePage";
type Props = {
    question: Question | null;
    scoreboard: Record<string, number> | null;
    view: string;
    playerId: string;
    hasAnswered: boolean;
    handleAnswer: (selected: string) => void;
};
const PlayerGame: React.FC<Props> = ({
    view,
    question,
    handleAnswer,
    hasAnswered,
    scoreboard,
    playerId,
}) => {
    return (
        <div>
            <div>player</div>
            <div>{view}</div>
            {view === "question" && question !== null && (
                <div>
                    <div>{question.trackName}</div>
                    <div>
                        {question.options.map((option) => {
                            return (
                                <button
                                    className="border"
                                    disabled={hasAnswered}
                                    onClick={() => handleAnswer(option)}
                                >
                                    {option}
                                </button>
                            );
                        })}
                    </div>
                </div>
            )}
            {view !== "question" && <div>Waiting for next question</div>}
        </div>
    );
};

export default PlayerGame;
