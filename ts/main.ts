import * as Nightmare from 'nightmare'
import { CPU, GPU, Component } from './component'

let nightmare = new Nightmare({show: false})

async function get() {
    try {
        let result = await nightmare
            .goto("http://cpu.userbenchmark.com/Intel-Core-i7-8700K/Rating/3937")
            .evaluate(() => { return document.querySelector('html')!.innerHTML })
            .end();
        return result
    } catch(e) {
        console.log(e)
        return undefined
    }
}

(async function() {
    console.log(await get())
}())