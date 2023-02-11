import http from 'k6/http';
import { check } from 'k6';
import { server } from './config.js';

export const options = {
    scenarios: {
        send: {
            executor: 'shared-iterations',
            vus: 1, // number of threads
            iterations: server.PerVUsIter,
            maxDuration: '60s',
        },
    },
};


const TotalCount = server.VUs * server.PerVUsIter

function Check_Record(r) {
    if (r.length != server.VUs) return false
    for (let t = 0; t < r.length; ++t) {
        let res = r[t].location == "l" + cnt && r[t].timestamp == server.ts + "T00:00:00.000+08:00"
        res &= r[t].signature == "NA==" && r[t].material == 19
        res &= (r[t].a == 1 && r[t].b == 1 && r[t].c == 1 && r[t].d == 1)
        if (!res) return false;
    }
    return true
}


function Check_Report(r) {
    let res = r.count == server.VUs && r.material == server.VUs * 19
    res &= (r.a == server.VUs && r.b == server.VUs && r.c == server.VUs && r.d == server.VUs)
    return res
}

let cnt = 0

export default function () {
    cnt++
    // const res1 = http.get(`https://www.google.com.tw?location=l${cnt}&date=2023-01-01`)
    // const res2 = http.get(`https://www.google.com.tw?location=l${cnt}&date=2023-01-01`)
    const res1 = http.get(`http://${server.host}:${server.port}/api/record?location=l${cnt}&date=${server.ts}`)
    const res2 = http.get(`http://${server.host}:${server.port}/api/report?location=l${cnt}&date=${server.ts}`)
    console.log(res1)
    check(res1, {
        '> Record Status is 200': (r) => r.status === 200,
        '> Record is Correct': (r) => Check_Record(JSON.parse(r.body)),
    });
    check(res2, {
        '> Report Status is 200': (r) => r.status === 200,
        '> Report is Correct': (r) => Check_Report(JSON.parse(r.body)),
    });
}
