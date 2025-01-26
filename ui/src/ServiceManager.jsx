import React, { useState } from 'react';

function ServiceManager() {
    const [services, setServices] = useState([
        { name: "personalSite", status: "stopped" },
        { name: "basicCalculator", status: "stopped" }
    ]);
    const [loading, setLoading] = useState(false);
    const [loadingService, setLoadingService] = useState(null); // Track which service is loading

    const backendUrl = 'http://localhost:8080';

    const handleStart = async (serviceName) => {
        setLoading(true);
        setLoadingService(serviceName); // Set the service being loaded
        const res = await fetch(`${backendUrl}/api/start/${serviceName}`, { method: 'POST' });
        setServices(services.map(service =>
            service.name === serviceName
                ? { ...service, status: "running" }
                : service
        ));
        setLoading(false);
        setLoadingService(null); // Reset loading state
        console.log('performed start');
    };

    const handleStop = async (serviceName) => {
        setLoading(true);
        setLoadingService(serviceName);
        await fetch(`${backendUrl}/api/stop/${serviceName}`, { method: 'POST' });
        setServices(services.map(service =>
            service.name === serviceName
                ? { ...service, status: "stopped" }
                : service
        ));
        setLoading(false);
        setLoadingService(null);
        console.log('performed stop');
    };

    const handleRestart = async (serviceName) => {
        setLoading(true);
        setLoadingService(serviceName);
        await fetch(`${backendUrl}/api/restart/${serviceName}`, { method: 'POST' });
        setServices(services.map(service =>
            service.name === serviceName
                ? { ...service, status: "running" }
                : service
        ));
        setLoading(false);
        setLoadingService(null);
        console.log('performed restart');
    };

    return (
        <div>
            <h1>Service Manager</h1>
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
                            <td>{service.status}</td>
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
