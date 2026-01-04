import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 10000 }, // ramp up to 10000 connections
        { duration: '1m', target: 10000 },  // stay at 10000 connections
        { duration: '30s', target: 0 },  // ramp down
    ],
};

const BASE_URL = __ENV.WS_URL || 'ws://localhost:8080/ws/time';

export default function () {
    const url = BASE_URL;
    const params = { tags: { my_tag: 'hello' } };

    const res = ws.connect(url, params, function (socket) {
        socket.on('open', () => {
            console.log('connected');
            socket.send(JSON.stringify({
                action: 'subscribe',
                timezone: 'Asia/Kuala_Lumpur',
                format: '24hour'
            }));
        });

        socket.on('message', (data) => {
            const msg = JSON.parse(data);
            check(msg, {
                'type is time_update': (m) => m.type === 'time_update',
            });
        });

        socket.on('close', () => console.log('disconnected'));

        // Stay connected for 10 seconds to receive ticks
        socket.setTimeout(() => {
            socket.close();
        }, 10000);
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}
