import { Alert, AlertDescription, AlertTitle } from "./ui/alert";
import { AlertCircleIcon } from "lucide-react";

type Props = {
    title: string | number | null | undefined;
    msg: string | null;
};
const CustomAlert: React.FC<Props> = ({ msg, title }) => {
    return (
        <div className="fixed top-24 left-1/2 -translate-x-1/2 z-[999] w-full max-w-md px-4">
            <Alert
                variant="destructive"
                className="bg-white text-red-600 border border-red-300 shadow-lg animate-in slide-in-from-top duration-500 ease-out"
            >
                <AlertCircleIcon className="h-5 w-5 text-red-500" />
                <AlertTitle className="text-red-700 font-semibold">{title}</AlertTitle>
                <AlertDescription className="text-sm text-red-500">
                    {msg}
                </AlertDescription>
            </Alert>
        </div>
    );
};

export default CustomAlert;
