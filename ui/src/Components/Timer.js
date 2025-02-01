import React, { useEffect, useState } from "react";

const Timer = ({ customTime }) => {
    const timeLimit = customTime / 1000 || 5 // Set default as 5 seconds
    const [countdown, setCountdown] = useState(timeLimit)
    const [lastRefreshed, setLastRefreshed] = useState(new Date())

    useEffect(() => {
        const interval = setInterval(() => {
            setCountdown((prev) => {
                if (prev === 1) {
                    return timeLimit
                }
                return prev - 1
            })
        }, 1000)
        setLastRefreshed(new Date())

        return () => clearInterval(interval)
    }, [timeLimit]);

    return (
        <div style={{ fontSize: '16px' }}>
            <p>Time until next refresh: <span style={{ color: 'yellow' }}>{countdown}</span></p>
            <p>Last refreshed: <span style={{ color: 'yellow' }}>
                {lastRefreshed.toLocaleString('en-US', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                    hour: 'numeric',
                    minute: '2-digit',
                    hour12: true
                })}</span></p>
        </div>
    );
};

export default Timer;
