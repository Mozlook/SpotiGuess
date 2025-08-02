import { useEffect, useState } from "react";
import type { Question } from "../pages/GamePage";

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
    const [position, setPosition] = useState<number>(1);
    const name = localStorage.getItem("name");
    useEffect(() => {
        const pos = scoreboard
            ? Object.entries(scoreboard)
                .sort(([, a], [, b]) => b - a)
                .findIndex(([id]) => id === name) + 1
            : 1;

        setPosition(pos);
    }, [scoreboard, playerId]);
    return (
        <div className="w-full max-w-2xl bg-white text-gray-800 p-6 rounded-xl shadow-lg flex flex-col items-center border border-gray-200">
            <div className="text-sm text-indigo-500 font-medium mb-3 uppercase tracking-wide">
                Player Mode
            </div>

            <div className="text-xl font-semibold mb-6 text-center">
                {view === "question"
                    ? "Answer the question!"
                    : "Waiting for next round..."}
            </div>

            {view === "question" && question && (
                <div className="w-full">
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        {question.options.map((option) => (
                            <button
                                key={option}
                                disabled={hasAnswered}
                                onClick={() => handleAnswer(option)}
                                className={`w-full px-4 py-3 rounded-lg text-left text-sm font-medium transition duration-200 border ${hasAnswered
                                        ? "bg-gray-100 text-gray-400 border-gray-300 cursor-not-allowed"
                                        : "bg-indigo-50 hover:bg-indigo-100 text-indigo-800 border-indigo-300"
                                    }`}
                            >
                                {option}
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {view !== "question" && (
                <div className="text-center text-sm text-gray-500 mt-6">
                    Please wait for the next round to begin.
                </div>
            )}

            <div className="mt-6 text-sm text-gray-500">
                Your position:{" "}
                <span className="font-semibold text-indigo-600">
                    {position > 0 ? `#${position}` : "?"}
                </span>
            </div>
        </div>
    );
};

export default PlayerGame;
