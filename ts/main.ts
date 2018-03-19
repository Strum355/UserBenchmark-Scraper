import * as Nightmare from 'nightmare'
import * as fs from 'fs'
import { JSDOM } from 'jsdom'
import { CPU, GPU, RAM, SSD, HDD, USB, Component } from './component'

let nightmare = new Nightmare({show: true})

interface user {
    username: string
    password: string
}

let u: user = JSON.parse(fs.readFileSync('./conf.json', 'utf-8'))

let login = async (nightmare: Nightmare) => {
    try {
        await nightmare
            .useragent('Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36')
            .goto('http://www.userbenchmark.com/page/login')
            .wait('input[name="username"]')
            .insert('input[name="username"]', u.username)
            .wait(1000)
            .insert('input[name="password"]', u.password)
            .wait(1000)
            .click('button[name="submit"]')
    } catch(e) {
        console.error(e)
    }
}

let scrape = async (nightmare: Nightmare, comp: Component) => {
    try {
        await nightmare
        .goto(comp.url)
        .wait('body')
        .evaluate((): string => {
            return document.querySelector('html')!.innerHTML
        })
        .then((res: string) => {
            const { document } = (new JSDOM(res)).window
            let body = document.querySelector('body')
            comp.get(body!, nightmare)
        })
    } catch(e) {
        console.error(e)
    }
}

async function run(nightmare: Nightmare) {
    try {
        let c = new CPU()
        c.url = "http://cpu.userbenchmark.com/Intel-Core-i7-8700K/Rating/3937"

        await login(nightmare)
        await scrape(nightmare, c)
        console.log(c)
    } catch(e) {
        console.error(e)
    }
}

run(nightmare)
    .then(nightmare.end)
    .catch(console.log)