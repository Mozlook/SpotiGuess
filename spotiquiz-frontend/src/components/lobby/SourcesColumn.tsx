import React from "react";
import type { GameMode, SearchItem } from "../../types/lobby";

type Props = {
    gameMode: GameMode;
    searchQuery: string;
    searchResults: SearchItem[];
    selectedSources: SearchItem[];
    onChangeGameMode: (mode: GameMode) => void;
    onChangeSearchQuery: (q: string) => void;
    onSearch: () => void;
    onSelectSource: (item: SearchItem) => void;
    onRemoveSource: (id: string) => void;
};

const SourcesColumn: React.FC<Props> = ({
    gameMode,
    searchQuery,
    searchResults,
    selectedSources,
    onChangeGameMode,
    onChangeSearchQuery,
    onSearch,
    onSelectSource,
    onRemoveSource,
}) => {
    return (
        <section className="bg-white/70 backdrop-blur border border-gray-200 rounded-xl shadow-md p-5 flex flex-col gap-4 overflow-auto">
            <h2 className="text-xl font-semibold text-gray-800">Configure Sources</h2>

            <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-gray-600">
                    Select Game Mode
                </label>
                <select
                    value={gameMode}
                    onChange={(e) => onChangeGameMode(e.target.value as GameMode)}
                    className="p-2 rounded border bg-white text-gray-800"
                >
                    <option value="players">Based on players</option>
                    <option value="playlist">From a playlist</option>
                    <option value="artist">From an artist</option>
                </select>
            </div>

            {(gameMode === "playlist" || gameMode === "artist") && (
                <>
                    {selectedSources.length > 0 && (
                        <div className="flex flex-wrap gap-3">
                            {selectedSources.map((source) => (
                                <div
                                    key={source.id}
                                    className="inline-flex items-center gap-3 bg-gray-500 text-white px-4 py-2 rounded-lg shadow-md"
                                >
                                    <img
                                        src={
                                            source.image ||
                                            "https://firstbenefits.org/wp-content/uploads/2017/10/placeholder-1024x1024.png"
                                        }
                                        alt={source.name}
                                        className="w-8 h-8 object-cover rounded"
                                    />
                                    <span className="text-base font-medium">{source.name}</span>
                                    <button
                                        onClick={() => onRemoveSource(source.id)}
                                        className="ml-2 text-white hover:text-red-400 font-bold text-lg"
                                        aria-label={`Remove ${source.name}`}
                                        title="Remove"
                                    >
                                        ×
                                    </button>
                                </div>
                            ))}
                        </div>
                    )}

                    <div className="flex gap-2">
                        <input
                            type="text"
                            placeholder={`Search ${gameMode}`}
                            value={searchQuery}
                            onChange={(e) => onChangeSearchQuery(e.target.value)}
                            onKeyDown={(e) => e.key === "Enter" && onSearch()}
                            className="px-4 py-2 rounded border bg-white text-gray-800 flex-1"
                        />
                        <button
                            onClick={onSearch}
                            className="bg-indigo-500 hover:bg-indigo-600 text-white px-4 py-2 rounded shadow text-sm"
                        >
                            Search
                        </button>
                    </div>

                    {searchResults.length > 0 && (
                        <div className="w-full bg-white border border-gray-200 rounded-lg shadow-md p-2 flex flex-col gap-2 max-h-72 overflow-auto">
                            {searchResults.map((result) => (
                                <button
                                    key={result.id}
                                    onClick={() => onSelectSource(result)}
                                    className="flex items-center gap-3 p-2 rounded hover:bg-gray-50 transition text-left"
                                >
                                    <img
                                        src={
                                            result.image ||
                                            "https://firstbenefits.org/wp-content/uploads/2017/10/placeholder-1024x1024.png"
                                        }
                                        alt={result.name}
                                        className="w-10 h-10 object-cover rounded"
                                    />
                                    <span className="text-sm font-medium text-gray-700">
                                        {result.name}
                                    </span>
                                </button>
                            ))}
                        </div>
                    )}
                </>
            )}
        </section>
    );
};

export default SourcesColumn;
