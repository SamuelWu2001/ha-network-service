import http from 'k6/http';
import { check } from 'k6';
import { server } from './config.js';

export const options = {
    scenarios: {
        send: {
            executor: 'per-vu-iterations',
            vus: server.VUs, // number of threads
            iterations: server.PerVUsIter,
            maxDuration: '60s',
        },
    },
};


let data = { "location": "", "timestamp": server.ts + "T00:00:00.000+08:00", "data": { "a": 1, "b": 1, "c": 1, "d": 1 } }
let cnt = 1;

export default function () {
    data.location = 'l' + cnt++
    // const res = http.post(`https://www.google.com.tw/`, JSON.stringify(data), server.param)
    const res = http.post(`http://${server.host}:${server.port}/api/order`, JSON.stringify(data), server.param)
    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
