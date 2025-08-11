import { useEffect, useState } from "react";
import axios from "axios";
import type { Question } from "../pages/GamePage";
import TimedProgress from "./TimedProgress";

type Props = {
    question: Question | null;
    scoreboard: Record<string, number> | null;
    view: string;
    hasAnswered: boolean;
    setHasAnswered: React.Dispatch<React.SetStateAction<boolean>>;
    code: string | undefined;
    playerName: string;
};

const PlayerGame: React.FC<Props> = ({
    view,
    question,
    hasAnswered,
    setHasAnswered,
    scoreboard,
    code,
    playerName,
}) => {
    const [position, setPosition] = useState<number>(1);
    const [selectedAnswer, setSelectedAnswer] = useState<string | null>(null);
    const [isCorrect, setIsCorrect] = useState<boolean | null>(null);
    const [earnedPoints, setEarnedPoints] = useState<number | null>();
    const apiUrl = import.meta.env.VITE_BACKEND.API.URL;
    const name = localStorage.getItem("name");
    useEffect(() => {
        const pos = scoreboard
            ? Object.entries(scoreboard)
                .sort(([, a], [, b]) => b - a)
                .findIndex(([id]) => id === name) + 1
            : 1;

        setPosition(pos);
    }, [scoreboard, name]);
    async function handleAnswer(selected: string) {
        if (!question) return;

        try {
            const res = await axios.post(`${apiUrl}/submit-answer`, {
                roomCode: code,
                questionId: question.id,
                selected,
                playerID: playerName,
            });
            setSelectedAnswer(selected);
            setIsCorrect(res.data.correct);
            setEarnedPoints(res.data.earned);

            setHasAnswered(true);
        } catch (err) {
            console.error("Submit failed:", err);
        }
    }
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
                <>
                    <div className="w-full mb-4">
                        <TimedProgress duration={15} />
                    </div>
                    <div className="w-full">
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            {question.options.map((option) => {
                                let buttonStyle =
                                    "bg-indigo-50 hover:bg-indigo-100 text-indigo-800 border-indigo-300";

                                if (hasAnswered) {
                                    if (
                                        option === question.correct &&
                                        selectedAnswer !== option
                                    ) {
                                        buttonStyle =
                                            "bg-green-100 text-green-700 border-green-300";
                                    } else if (option === selectedAnswer && isCorrect) {
                                        buttonStyle = "bg-green-500 text-white border-green-600";
                                    } else if (option === selectedAnswer && !isCorrect) {
                                        buttonStyle = "bg-red-500 text-white border-red-600";
                                    } else {
                                        buttonStyle = "bg-gray-100 text-gray-400 border-gray-300";
                                    }
                                }

                                return (
                                    <button
                                        key={option}
                                        disabled={hasAnswered}
                                        onClick={() => handleAnswer(option)}
                                        className={`w-full px-4 py-3 rounded-lg text-left text-sm font-medium transition duration-200 border ${buttonStyle}`}
                                    >
                                        {option}
                                    </button>
                                );
                            })}
                        </div>
                        {hasAnswered && isCorrect && earnedPoints && (
                            <div className="mt-4 text-green-600 font-semibold text-center text-lg">
                                Correct! +{earnedPoints} points
                            </div>
                        )}
                    </div>
                </>
            )}

            {view !== "question" && (
                <>
                    <div className="w-full mb-4">
                        <TimedProgress duration={5} />
                    </div>
                    <div className="text-center text-sm text-gray-500">
                        Please wait for the next round to begin.
                    </div>
                </>
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
