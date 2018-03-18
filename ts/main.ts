import * as Nightmare from 'nightmare'
import * as fs from 'fs'
import { CPU, GPU, RAM, SSD, HDD, USB, Component } from './component'

let nightmare = new Nightmare({show: true})

interface user {
    username: string
    password: string
}

let u: user = JSON.parse(fs.readFileSync('./config.json', 'utf-8'))
console.log(u)
let c = new CPU()
c.url = "http://cpu.userbenchmark.com/Intel-Core-i7-8700K/Rating/3937"

try {
    login(nightmare)
/*     nightmare
        .goto(c.url)
        .wait('body')
        .evaluate((): string => {
            return document.querySelector('body')!.innerHTML
        })
        .end()
        .then(console.log)
        .catch((e) => {
            console.error(e)
        }) */
} catch(e) {
    console.error(e)
}

function login(nightmare: Nightmare){
    try {
        nightmare
            .goto('http://www.userbenchmark.com/page/login')
            .exists('input[name="username"]', (exists: boolean) => {
                if(exists) {
                    console.log("exists")
                } else {
                    console.log("doesnt exist")
                }
            })
            .type('input[name="username"]', "test")
            .wait(2000)
            .type('input[name="password"]', "test")
            .wait(4000)
            .evaluate((): string => {
                return (document.querySelector('input[name="username"]')! as HTMLInputElement).value
            })
            .then((res: string) => {
                console.log(res)
            })
            .catch((e) => {
                console.error(e)
            })
        nightmare
            .click('button[name="submit"]')
    } catch(e) {
        console.error(e)
    }
}
