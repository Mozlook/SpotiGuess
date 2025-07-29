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
                className="bg-blue-600 px-4 py-2 rounded text-white"
            >
                Join
            </DialogTrigger>
            <DialogContent className="bg-gray-900 text-white border-none">
                <DialogHeader>
                    <DialogTitle className="text-lg font-bold">
                        What's your name?
                    </DialogTitle>
                </DialogHeader>
                <input
                    type="text"
                    placeholder="Display name"
                    onChange={(e) => setName(e.target.value)}
                    className="px-4 py-2 w-full rounded bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-600"
                />
                <DialogFooter>
                    <DialogClose>
                        <button>Cancel</button>
                    </DialogClose>
                    <button onClick={() => onConfirm(name)}>Join</button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};
export default CustomDialog;
