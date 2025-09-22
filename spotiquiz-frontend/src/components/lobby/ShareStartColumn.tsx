import React from "react";
import QRPanel from "../QRPanel";

type Props = {
    roomCode: string;
    loading: boolean;
    errorMessage?: string;
    onLoadPrevious: () => void;
    onStartGame: () => void;
};

const ShareStartColumn: React.FC<Props> = ({
    roomCode,
    loading,
    errorMessage,
    onLoadPrevious,
    onStartGame,
}) => {
    return (
        <section className="bg-white/70 backdrop-blur border border-gray-200 rounded-xl shadow-md p-5 flex flex-col items-center gap-4 overflow-auto">
            <QRPanel roomCode={roomCode} />

            <p className="text-indigo-700 font-medium">You are the host</p>

            <div className="flex flex-col gap-2 w-full max-w-sm">
                <button
                    onClick={onLoadPrevious}
                    className="w-full bg-green-900 hover:bg-green-700 text-white font-medium py-2 px-4 rounded shadow text-sm"
                >
                    Load previous game
                </button>
                <button
                    onClick={onStartGame}
                    disabled={loading}
                    className={`w-full ${loading ? "bg-gray-500" : "bg-indigo-600 hover:bg-indigo-700"
                        } text-white font-medium py-2 px-4 rounded shadow text-sm`}
                >
                    {loading ? "Creating…" : "Start game"}
                </button>
            </div>

            {errorMessage && (
                <p className="text-red-600 font-medium text-center">{errorMessage}</p>
            )}
        </section>
    );
};

export default ShareStartColumn;
