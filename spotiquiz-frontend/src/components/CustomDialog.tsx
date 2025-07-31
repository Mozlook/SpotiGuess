import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogClose,
} from "./ui/dialog";
import { useState } from "react";
type Props = {
    onConfirm: (name: string) => void;
    roomCode: string;
};
const CustomDialog: React.FC<Props> = ({ onConfirm, roomCode }) => {
    const [name, setName] = useState<string>("");

    return (
        <Dialog>
            <DialogTrigger
                disabled={roomCode.trim().length !== 6}
                className={`px-4 py-2 rounded font-semibold shadow transition ${roomCode.trim().length === 6
                        ? "bg-indigo-600 text-white hover:bg-indigo-700"
                        : "bg-gray-300 text-gray-500 cursor-not-allowed"
                    }`}
            >
                Join
            </DialogTrigger>

            <DialogContent className="bg-gray-300 text-gray-800 border border-gray-200 rounded-lg max-w-sm">
                <DialogHeader>
                    <DialogTitle className="text-lg font-bold text-center">
                        What's your name?
                    </DialogTitle>
                </DialogHeader>

                <input
                    type="text"
                    placeholder="Display name"
                    onChange={(e) => setName(e.target.value)}
                    className="w-full px-4 py-2 rounded bg-gray-100 text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />

                <DialogFooter className="flex justify-end gap-2 mt-4">
                    <DialogClose>
                        <button className="px-4 py-2 rounded bg-rose-500 text-white hover:bg-rose-600 transition">
                            Cancel
                        </button>
                    </DialogClose>
                    <button
                        onClick={() => onConfirm(name)}
                        className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded transition"
                    >
                        Join
                    </button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};
export default CustomDialog;
