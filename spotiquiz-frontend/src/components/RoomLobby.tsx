import { useParams } from "react-router-dom";
const RoomLobby = () => {
    const { code } = useParams();
    const isHost = localStorage.getItem("isHost");

    return (
        <>
            <span>Room code:{code}</span>
            {isHost ? (
                <div>
                    <span>You are a host</span>
                    <button>start gry</button>
                </div>
            ) : (
                <span>Waiting for host</span>
            )}
        </>
    );
};
export default RoomLobby;
