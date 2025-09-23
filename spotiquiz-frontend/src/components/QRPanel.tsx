type QRPanelProps = { roomCode: string };
import { QRCodeCanvas } from "qrcode.react";

const QRPanel: React.FC<QRPanelProps> = ({ roomCode }) => {
    const joinUrl = `https://spotiguess.mmozoluk.com/${roomCode}`;

    const canvasId = `qr-${roomCode}`;

    return (
        <div className="w-full bg-white/70 backdrop-blur border border-gray-200 rounded-xl shadow-md p-5 flex flex-col items-center gap-3">
            <h2 className="text-xl font-semibold text-gray-800">Scan to join</h2>

            <div className="bg-white p-2 rounded-xl shadow-inner">
                <QRCodeCanvas id={canvasId} value={joinUrl} size={224} level="M" />
            </div>
        </div>
    );
};
export default QRPanel;
