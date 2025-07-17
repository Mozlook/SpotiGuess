import type { Question } from "./GamePage";
type Props = {
    question: Question | null;
    scoreboard: Record<string, number> | null;
    view: string;
    playerId: string;
};
const PlayerGame: React.FC<Props> = ({ view, question }) => {
    return (
        <div>
            <div>player</div>
            <div>{view}</div>
            {view === "question" && <div>{question?.trackName}</div>}
            {view == "scoreboard" && <div>Waiting for next question</div>}
        </div>
    );
};

export default PlayerGame;
