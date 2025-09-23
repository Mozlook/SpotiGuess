import React from "react";

type Props = {
    players: string[];
};

const PlayersColumn: React.FC<Props> = ({ players }) => {
    return (
        <section className="bg-white/70 backdrop-blur border border-gray-200 rounded-xl shadow-md p-5 flex flex-col overflow-auto">
            <h2 className="text-xl font-semibold text-gray-800 text-center mb-4">
                Players in Room
            </h2>
            {players.length > 0 ? (
                <ul className="space-y-2 pr-1">
                    {players.map((player, index) => (
                        <li
                            key={`${player}-${index}`}
                            className="flex items-center justify-between bg-gray-100 px-4 py-2 rounded-md text-sm font-medium text-gray-800"
                        >
                            <span className="truncate">{player}</span>
                        </li>
                    ))}
                </ul>
            ) : (
                <p className="text-gray-600 text-center">
                    No players yet. Share the QR or link.
                </p>
            )}
        </section>
    );
};

export default PlayersColumn;
