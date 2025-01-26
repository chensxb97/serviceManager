import React, { useState } from 'react';

function ServiceManager() {
    const [services, setServices] = useState([
        { name: "personalSite", status: "stopped" },
        { name: "basicCalculator", status: "stopped" }
    ]);

    const backendUrl = 'http://localhost:8080'

    const handleStart = async (serviceName) => {
        await fetch(`${backendUrl}/api/start/${serviceName}`, { method: 'POST' });

        setServices(services.map(service =>
            service.name === serviceName
                ? { ...service, status: "running" }
                : service
        ));
        console.log('performed start')
    };

    const handleStop = async (serviceName) => {
        await fetch(`${backendUrl}/api/stop/${serviceName}`, { method: 'POST' });

        setServices(services.map(service =>
            service.name === serviceName
                ? { ...service, status: "stopped" }
                : service
        ));
        console.log('performed stop')
    };

    const handleRestart = async (serviceName) => {
        await fetch(`${backendUrl}/api/restart/${serviceName}`, { method: 'POST' });

        setServices(services.map(service =>
            service.name === serviceName
                ? { ...service, status: "running" }
                : service
        ));
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
                                <button onClick={() => handleStart(service.name)}>Start</button>
                                <button onClick={() => handleStop(service.name)}>Stop</button>
                                <button onClick={() => handleRestart(service.name)}>Restart</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export default ServiceManager;