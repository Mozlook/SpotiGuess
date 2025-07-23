import { useLocation, useNavigate } from "react-router-dom";

type Scoreboard = Record<string, number>;

const ScoreboardPage = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const scoreboard = location.state as Scoreboard | null;
    console.log("Location:", location);
    console.log("State:", location.state);

    if (!scoreboard || Object.keys(scoreboard).length === 0) {
        return <p className="text-red-400">No scoreboard data available.</p>;
    }

    return (
        <div className=" bg-gray-900">
            <div className="min-h-screen max-w-md mx-auto text-white p-6">
                <button onClick={() => navigate("/")}>Home</button>
                <h1 className="text-2xl font-bold mb-4">Final Scoreboard</h1>
                <ul className="space-y-2">
                    {Object.entries(scoreboard)
                        .sort(([, a], [, b]) => b - a)
                        .map(([playerId, score], index) => (
                            <li
                                key={playerId}
                                className="flex justify-between bg-gray-800 p-2 rounded"
                            >
                                <span>
                                    {index + 1}. {playerId}
                                </span>
                                <span>{score} pts</span>
                            </li>
                        ))}
                </ul>
            </div>
        </div>
    );
};

export default ScoreboardPage;
