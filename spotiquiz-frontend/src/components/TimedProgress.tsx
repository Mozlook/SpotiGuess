import { Progress } from "./ui/progress";
import { useState, useEffect } from "react";
type Props = {
    duration: number;
};

const TimedProgress: React.FC<Props> = ({ duration }) => {
    const [progress, setProgress] = useState(0);
    const intervalMs = 50;
    const step = 100 / ((duration * 1000) / intervalMs);
    useEffect(() => {
        setProgress(0);
        const timer = setInterval(() => {
            setProgress((prev) => Math.min(prev + step, 100));
        }, intervalMs);

        return () => clearInterval(timer);
    }, [duration, step]);
    return (
        <Progress
            value={progress}
            className="h-3 w-full bg-gray-200 [&>*]:bg-indigo-500"
        />
    );
};

export default TimedProgress;
