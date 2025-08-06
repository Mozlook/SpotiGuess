import { useLocation, useNavigate } from "react-router-dom";

type Scoreboard = Record<string, number>;

const ScoreboardPage = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const scoreboard = location.state as Scoreboard | null;

    if (!scoreboard || Object.keys(scoreboard).length === 0) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100">
                <p className="text-red-600 text-lg font-medium">
                    No scoreboard data available.
                </p>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gradient-to-b from-emerald-300 via-gray-200 to-emerald-100 flex items-center justify-center px-4">
            <div className="w-full max-w-xl bg-white text-gray-800 p-6 rounded-xl shadow-lg border border-gray-200">
                <h1 className="text-3xl font-bold text-center text-indigo-700 mb-6">
                    Final Scoreboard
                </h1>

                <ul className="space-y-3">
                    {Object.entries(scoreboard)
                        .sort(([, a], [, b]) => b - a)
                        .map(([playerId, score], index) => (
                            <li
                                key={playerId}
                                className="flex justify-between items-center px-4 py-3 bg-indigo-50 border border-indigo-200 rounded-lg"
                            >
                                <span className="font-medium text-indigo-700">
                                    #{index + 1} {playerId}
                                </span>
                                <span className="text-indigo-600 font-semibold">
                                    {score} pts
                                </span>
                            </li>
                        ))}
                </ul>

                <div className="mt-6 text-center">
                    <button
                        onClick={() => navigate("/")}
                        className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded shadow transition"
                    >
                        Back to Home
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ScoreboardPage;
