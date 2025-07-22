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
}) => {
    return (
        <div className="w-full max-w-xl bg-gray-800 text-white p-6 rounded-xl shadow-md flex flex-col items-center">
            <div className="text-sm text-gray-400 mb-4">🎮 Player</div>
            <div className="text-lg font-semibold mb-6 capitalize">
                {view === "question"
                    ? "Answer the question!"
                    : "Waiting for next question..."}
            </div>

            {view === "question" && question && (
                <div className="w-full">
                    <h2 className="text-2xl font-bold text-center mb-4">
                        {question.trackName}
                    </h2>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        {question.options.map((option) => (
                            <button
                                key={option}
                                disabled={hasAnswered}
                                onClick={() => handleAnswer(option)}
                                className={`px-4 py-2 rounded-lg text-left border border-gray-600 transition duration-200 ${hasAnswered
                                        ? "bg-gray-700 text-gray-400 cursor-not-allowed"
                                        : "bg-gray-900 hover:bg-gray-700"
                                    }`}
                            >
                                {option}
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {view !== "question" && (
                <div className="text-center text-gray-400 mt-6 text-sm">
                    Please wait for the next round to begin.
                </div>
            )}
        </div>
    );
};

export default PlayerGame;
