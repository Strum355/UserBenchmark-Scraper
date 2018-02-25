import Nightmare = require("nightmare");

let nightmare = new Nightmare({show: true})

nightmare
    .goto("http://cpu.userbenchmark.com/Intel-Core-i7-8700K/Rating/3937")
    .end()
    .then(console.log)
    .catch(error => {
        console.error(error)
    })
