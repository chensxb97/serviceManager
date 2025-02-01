import React, { useState, useEffect } from 'react';
import Timer from './Components/Timer'
function ServiceManager() {
    const [animateKey, setAnimateKey] = useState(0)
    const timeLimit = 5000


    const [services, setServices] = useState([
        { name: "app1" },
        { name: "app2" }
    ]);

    const [loading, setLoading] = useState(false);
    const [loadingService, setLoadingService] = useState(null);

    const backendUrl = 'http://localhost:8080';

    const handleStart = async (serviceName) => {
        setLoading(true);
        setLoadingService(serviceName);
        try {
            const res = await fetch(`${backendUrl}/api/action`, {
                method: 'POST',
                body: JSON.stringify({
                    service: `${serviceName}`,
                    action: 'start'
                })
            });
            if (res?.status === 200) {
                alert('Start Operation Submitted Successfully')
                console.log('performed start');
            } else {
                console.log('error returned from start operation')
            }
        } catch (ex) {
            console.log('error while performing start operation')
            alert('Unable to perform start operation')
        }
        setLoading(false);
        setLoadingService(null);
    };

    const handleStop = async (serviceName) => {
        setLoading(true);
        setLoadingService(serviceName);
        try {
            const res = await fetch(`${backendUrl}/api/action`, {
                method: 'POST',
                body: JSON.stringify({
                    service: `${serviceName}`,
                    action: 'stop'
                })
            });
            if (res?.status === 200) {
                alert('Stop Operation Submitted Successfully')
                console.log('performed stop');
            } else {
                console.log('error in stop operation')
                alert('error in stop operation')
            }
        }
        catch (ex) {
            console.log('error while performing stop operation')
            alert('error while performing stop operation')
        }
        setLoading(false);
        setLoadingService(null);

    };

    const handleRestart = async (serviceName) => {
        setLoading(true);
        setLoadingService(serviceName);
        const res = await fetch(`${backendUrl}/api/action`, {
            method: 'POST',
            body: JSON.stringify({
                service: `${serviceName}`,
                action: 'restart'
            })
        });
        try {
            if (res?.status === 200) {
                alert('Restart Operation Submitted Successfully')
                console.log('performed restart');
            } else {
                console.log('error in restart operation')
                alert('error in restart operation')
            }
        }
        catch (ex) {
            console.log('error while performing restart operation')
            alert('error while performing restart operation')
        }
        setLoading(false);
        setLoadingService(null);
    };


    const fetchServiceStates = async () => {
        try {
            const res = await fetch(`${backendUrl}/api/services`, {
                method: 'GET'
            })
            if (res.status === 200) {
                const serviceStateMap = await res.json()
                setServices((prevServices) =>
                    prevServices.map(service => ({
                        ...service,
                        status: serviceStateMap[service.name] || service.status
                    })))
            }
        } catch (ex) {
            console.error('error while fetching service states')
        }
    }

    useEffect(() => {
        fetchServiceStates()
        const intervalId = setInterval(() => {
            fetchServiceStates()
            setAnimateKey((prevKey) => !prevKey)
        }, timeLimit)
        return () => clearInterval(intervalId)
        // eslint-disable-next-line
    }, [])

    return (
        <div>
            <h1>Service Manager</h1>
            <Timer customTime={timeLimit} />
            <table>
                <thead>
                    <tr>
                        <th>Service Name</th>
                        <th>Status</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {services.map(service => (
                        <tr key={service.name}>
                            <td>{service.name}</td>
                            <td
                                key={animateKey}
                                style={{
                                    padding: "15px",
                                    color: service?.status === "running" ? "lightgreen" : "salmon",
                                    animation: service?.status ? "bounceOut 1s ease" : "",
                                }}
                            >
                                {service?.status}
                            </td>
                            <td>
                                <button
                                    onClick={() => handleStart(service.name)}
                                    disabled={loading && loadingService === service.name}
                                >
                                    {loadingService === service.name ? '...' : 'Start'}
                                </button>
                                <button
                                    onClick={() => handleStop(service.name)}
                                    disabled={loading && loadingService === service.name}
                                >
                                    {loading && loadingService === service.name ? '...' : 'Stop'}
                                </button>
                                <button
                                    onClick={() => handleRestart(service.name)}
                                    disabled={loading && loadingService === service.name}
                                >
                                    {loading && loadingService === service.name ? '...' : 'Restart'}
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default ServiceManager;
